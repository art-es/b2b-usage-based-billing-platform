package register

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/register/dto"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
)

func TestEndpoint(t *testing.T) {
	type deps struct {
		mockUsecase *Mockusecase
		logger      log.Logger
		logbuf      log.Buffer
	}

	newDeps := func() *deps {
		mockCtrl := gomock.NewController(t)
		logbuf := log.NewBuffer()

		return &deps{
			mockUsecase: NewMockusecase(mockCtrl),
			logger: log.NewLogger(&log.Options{
				Output:       logbuf,
				GetCreatedAt: func() string { return "test-created-at" },
			}),
			logbuf: logbuf,
		}
	}

	callEndpoint := func(d *deps, reqBody io.Reader) *httptest.ResponseRecorder {
		mux := http.NewServeMux()
		Bind(mux, d.mockUsecase, d.logger)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/v1/auth/register", reqBody)
		mux.ServeHTTP(w, r)

		return w
	}

	t.Run("registration accepted", func(t *testing.T) {
		d := newDeps()

		expUcReq := &dto.Request{
			Name:     "test-name",
			Email:    "test-email@example.com",
			Password: "test-password123",
		}

		d.mockUsecase.EXPECT().
			Do(gomock.Any(), gomock.Eq(expUcReq)).
			Return(nil)

		reqBody := strings.NewReader(`{
			"name": "test-name",
			"email": "test-email@example.com",
			"password": "test-password123"
		}`)

		resp := callEndpoint(d, reqBody)

		assert.Equal(t, 202, resp.Code)
		assert.Empty(t, resp.Body.String())
		assert.Empty(t, d.logbuf.Logs())
	})

	t.Run("unexpected usecase error", func(t *testing.T) {
		d := newDeps()

		expUcReq := &dto.Request{
			Name:     "test-name",
			Email:    "test-email@example.com",
			Password: "test-password123",
		}

		d.mockUsecase.EXPECT().
			Do(gomock.Any(), gomock.Eq(expUcReq)).
			Return(errors.New("test error"))

		reqBody := strings.NewReader(`{
			"name": "test-name",
			"email": "test-email@example.com",
			"password": "test-password123"
		}`)

		resp := callEndpoint(d, reqBody)

		assert.Equal(t, 500, resp.Code)
		assert.JSONEq(t, `{"message":"Internal error"}`, resp.Body.String())

		logs := d.logbuf.Logs()
		assert.Len(t, logs, 1)
		assert.JSONEq(t, `{
			"created_at": "test-created-at",
			"level": "error",
			"pkg": "internal/transport/http/endpoints/register",
			"message": "unexpected usecase error",
			"error": "test error"
		}`, logs[0])
	})

	t.Run("email is already in use", func(t *testing.T) {
		d := newDeps()

		expUcReq := &dto.Request{
			Name:     "test-name",
			Email:    "test-email@example.com",
			Password: "test-password123",
		}

		d.mockUsecase.EXPECT().
			Do(gomock.Any(), gomock.Eq(expUcReq)).
			Return(dto.ErrEmailInUse)

		reqBody := strings.NewReader(`{
			"name": "test-name",
			"email": "test-email@example.com",
			"password": "test-password123"
		}`)

		resp := callEndpoint(d, reqBody)

		assert.Equal(t, 400, resp.Code)
		assert.JSONEq(t, `{
			"code": 2007,
			"message": "Email is already in use"
		}`, resp.Body.String())
		assert.Empty(t, d.logbuf.Logs())
	})

	t.Run("request body validation error", func(t *testing.T) {
		d := newDeps()

		resp := callEndpoint(d, strings.NewReader(`{}`))

		assert.Equal(t, 400, resp.Code)
		assert.JSONEq(t, `{
			"code": 2001,
			"message": "Required name"
		}`, resp.Body.String())
		assert.Empty(t, d.logbuf.Logs())
	})

	t.Run("invalid request format error", func(t *testing.T) {
		d := newDeps()

		resp := callEndpoint(d, strings.NewReader(`foo`))

		assert.Equal(t, 400, resp.Code)
		assert.JSONEq(t, `{
			"code": 1001,
			"message": "Invalid request format"
		}`, resp.Body.String())
		assert.Empty(t, d.logbuf.Logs())
	})
}

func TestEndpoint_WrongHTTPMethod(t *testing.T) {
	logbuf := log.NewBuffer()
	mux := http.NewServeMux()

	Bind(
		mux,
		NewMockusecase(gomock.NewController(t)),
		log.NewLogger(&log.Options{Output: logbuf}),
	)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/v1/auth/register", nil)
	mux.ServeHTTP(w, r)

	assert.Equal(t, 405, w.Code)
	assert.Equal(t, "Method Not Allowed\n", w.Body.String())
	assert.Empty(t, logbuf.Logs())
}
