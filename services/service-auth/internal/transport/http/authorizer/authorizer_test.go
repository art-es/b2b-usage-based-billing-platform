//go:generate mockgen -source=authorizer.go -destination=authorizer_mock_test.go -package=$GOPACKAGE
package authorizer

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthorizer(t *testing.T) {
	testJwtSecret := "test_jwt_secret"
	testToken := "test_token"

	type deps struct {
		mockJwtService *MockjwtService
		logbuf         log.Buffer
		authorizer     *Authorizer
	}

	newDeps := func(t *testing.T) *deps {
		mockCtrl := gomock.NewController(t)
		mockJwtService := NewMockjwtService(mockCtrl)

		logbuf := log.NewBuffer()
		logger := log.NewLogger(&log.Options{
			Output:       logbuf,
			GetCreatedAt: func() string { return "test_created_at" },
		})

		return &deps{
			mockJwtService: mockJwtService,
			logbuf:         logbuf,
			authorizer:     New(mockJwtService, testJwtSecret, logger),
		}
	}

	run := func(d *deps, header string) (*httptest.ResponseRecorder, *jwt.Claims, bool) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("Authorization", header)
		claims, ok := d.authorizer.Authorize(w, r)
		return w, claims, ok
	}

	t.Run("ok", func(t *testing.T) {
		d := newDeps(t)

		parseResult := &jwt.Claims{SessionID: "test_session_id"}

		d.mockJwtService.EXPECT().
			Parse(gomock.Eq([]byte(testJwtSecret)), gomock.Eq(testToken)).
			Return(parseResult, nil)

		w, claims, ok := run(d, "Bearer "+testToken)

		assert.True(t, ok)
		assert.Equal(t, parseResult, claims)
		assert.Equal(t, 200, w.Code)
		assert.Empty(t, d.logbuf.Logs())
	})

	t.Run("unexpected jwt service parse error", func(t *testing.T) {
		d := newDeps(t)

		d.mockJwtService.EXPECT().
			Parse(gomock.Eq([]byte(testJwtSecret)), gomock.Eq(testToken)).
			Return(nil, errors.New("test error"))

		w, claims, ok := run(d, "Bearer "+testToken)

		assert.False(t, ok)
		assert.Nil(t, claims)
		assert.Equal(t, 500, w.Code)
		assert.JSONEq(t, `{"message":"Internal error"}`, w.Body.String())

		logs := d.logbuf.Logs()
		assert.Len(t, logs, 1)
		assert.JSONEq(t, `{
			"level": "error",
			"created_at": "test_created_at",
			"pkg": "internal/http/authorizer",
			"message": "unexpected jwt service parse error",
			"error": "test error"
		}`, logs[0])
	})

	t.Run("invalid token", func(t *testing.T) {
		d := newDeps(t)

		d.mockJwtService.EXPECT().
			Parse(gomock.Eq([]byte(testJwtSecret)), gomock.Eq(testToken)).
			Return(nil, jwt.ErrInvalidToken)

		w, claims, ok := run(d, "Bearer "+testToken)

		assert.False(t, ok)
		assert.Nil(t, claims)
		assertUnauthorized(t, w)
		assert.Empty(t, d.logbuf.Logs())
	})

	t.Run("invalid header format #1", func(t *testing.T) {
		d := newDeps(t)

		w, claims, ok := run(d, "Foo "+testToken)

		assert.False(t, ok)
		assert.Nil(t, claims)
		assertUnauthorized(t, w)
		assert.Empty(t, d.logbuf.Logs())
	})

	t.Run("invalid header format #2", func(t *testing.T) {
		d := newDeps(t)

		w, claims, ok := run(d, "Bearer "+testToken+" Foo")

		assert.False(t, ok)
		assert.Nil(t, claims)
		assertUnauthorized(t, w)
		assert.Empty(t, d.logbuf.Logs())
	})

	t.Run("no header", func(t *testing.T) {
		d := newDeps(t)

		w, claims, ok := run(d, "")

		assert.False(t, ok)
		assert.Nil(t, claims)
		assertUnauthorized(t, w)
		assert.Empty(t, d.logbuf.Logs())
	})
}

func assertUnauthorized(t *testing.T, w *httptest.ResponseRecorder) {
	assert.Equal(t, 401, w.Code)
	assert.JSONEq(t, `{"message":"Unauthorized"}`, w.Body.String())
}
