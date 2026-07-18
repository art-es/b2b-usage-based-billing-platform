package jwt

import (
	"testing"
	"time"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
	domain "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_OK(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockTimeService := NewMocktimeService(mockCtrl)
	mockTimeService.EXPECT().
		GetCurrentTime().
		Return(getTime("2026-01-02 12:00:00"))

	logbuf := log.NewBuffer()
	logger := log.NewLogger(&log.Options{Output: logbuf})

	svc := NewService(mockTimeService, logger)

	claims := domain.Claims{
		SessionID:      "test-session-id",
		UserID:         "test-user-id",
		OrgnID: ptr.To("test-org-id"),
		ExpiresAt:      getTime("2026-01-03 12:00:00"),
	}

	gotToken, err := svc.Generate([]byte("test-secret"), &claims)
	require.NoError(t, err)
	assert.NotEmpty(t, gotToken)

	gotClaims, err := svc.Parse([]byte("test-secret"), gotToken)
	require.NoError(t, err)
	require.NotNil(t, gotClaims)

	gotClaims.ExpiresAt = gotClaims.ExpiresAt.UTC()
	assert.Equal(t, claims, *gotClaims)
}

func TestService_Expired(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockTimeService := NewMocktimeService(mockCtrl)
	mockTimeService.EXPECT().
		GetCurrentTime().
		Return(getTime("2026-01-02 12:00:00"))

	logbuf := log.NewBuffer()
	logger := log.NewLogger(&log.Options{Output: logbuf})

	svc := NewService(mockTimeService, logger)

	claims := domain.Claims{
		SessionID:      "test-session-id",
		UserID:         "test-user-id",
		OrgnID: ptr.To("test-org-id"),
		ExpiresAt:      getTime("2026-01-01 12:00:00"),
	}

	gotToken, err := svc.Generate([]byte("test-secret"), &claims)
	require.NoError(t, err)
	assert.NotEmpty(t, gotToken)

	gotClaims, err := svc.Parse([]byte("test-secret"), gotToken)
	require.ErrorIs(t, err, jwt.ErrInvalidToken)
	require.Nil(t, gotClaims)
}

func getTime(s string) time.Time {
	out, _ := time.Parse(time.DateTime, s)
	return out.UTC()
}
