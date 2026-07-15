package refresh_session

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/session"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/repository"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/refresh_session/dto"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/trx"
)

func TestUsecase(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 7, 15, 12, 0, 0, 0, time.UTC)

	type deps struct {
		mockJWTService        *MockjwtService
		mockKeyedHashService  *MockkeyedHashService
		mockTimeService       *MocktimeService
		mockUUIDService       *MockuuidService
		mockSessionRepository *MocksessionRepository
		logbuf                log.Buffer
		usecase               *Usecase
	}

	newDeps := func() *deps {
		mockCtrl := gomock.NewController(t)
		mockJWTService := NewMockjwtService(mockCtrl)
		mockKeyedHashService := NewMockkeyedHashService(mockCtrl)
		mockTimeService := NewMocktimeService(mockCtrl)
		mockUUIDService := NewMockuuidService(mockCtrl)
		mockSessionRepository := NewMocksessionRepository(mockCtrl)

		logbuf := log.NewBuffer()
		logger := log.NewLogger(&log.Options{
			Output:       logbuf,
			GetCreatedAt: func() string { return "test-created-at" },
		})

		return &deps{
			mockJWTService:        mockJWTService,
			mockKeyedHashService:  mockKeyedHashService,
			mockTimeService:       mockTimeService,
			mockUUIDService:       mockUUIDService,
			mockSessionRepository: mockSessionRepository,
			logbuf:                logbuf,
			usecase: NewUsecase(
				mockJWTService,
				mockKeyedHashService,
				mockTimeService,
				mockUUIDService,
				mockSessionRepository,
				"test-jwt-secret",
				"test-refresh-secret",
				logger,
			),
		}
	}

	t.Run("ok", func(t *testing.T) {
		d := newDeps()

		orgID := "test-organization-id"
		ses := &session.Session{
			ID:               "test-session-id",
			UserID:           "test-user-id",
			OrganizationID:   &orgID,
			RefreshTokenHash: "old-refresh-token-hash",
		}

		d.mockTimeService.EXPECT().GetCurrentTime().Return(now)
		gomock.InOrder(
			d.mockKeyedHashService.EXPECT().
				Generate(gomock.Eq([]byte("test-refresh-secret")), gomock.Eq("old-refresh-token")).
				Return("old-refresh-token-hash", nil),
			d.mockKeyedHashService.EXPECT().
				Generate(gomock.Eq([]byte("test-refresh-secret")), gomock.Eq("new-refresh-token")).
				Return("new-refresh-token-hash", nil),
		)

		d.mockSessionRepository.EXPECT().
			GetByRefreshTokenHash(gomock.Any(), gomock.Eq("old-refresh-token-hash")).
			Return(ses, nil)

		d.mockUUIDService.EXPECT().
			Generate().
			Return("new-refresh-token")

		expSession := &session.Session{
			ID:                    "test-session-id",
			UserID:                "test-user-id",
			OrganizationID:        &orgID,
			RefreshTokenHash:      "new-refresh-token-hash",
			RefreshTokenExpiresAt: now.Add(session.RefreshTokenExpiry),
		}
		d.mockSessionRepository.EXPECT().
			Save(gomock.Any(), gomock.Eq(expSession)).
			Return(nil)

		d.mockJWTService.EXPECT().
			Generate(gomock.Eq([]byte("test-jwt-secret")), gomock.Eq(jwt.NewClaims("test-session-id", "test-user-id", &orgID))).
			Return("test-access-token", nil)

		res, err := d.usecase.Do(ctx, "old-refresh-token")
		assert.NoError(t, err)
		assert.Equal(t, &dto.Response{
			AccessToken:  "test-access-token",
			RefreshToken: "new-refresh-token",
		}, res)
		assert.Empty(t, d.logbuf.Logs())
	})

	t.Run("commit trx", func(t *testing.T) {
		d := newDeps()

		ses := &session.Session{ID: "test-session-id", UserID: "test-user-id"}

		d.mockTimeService.EXPECT().GetCurrentTime().Return(now)
		d.mockKeyedHashService.EXPECT().Generate(gomock.Any(), gomock.Any()).Return("old-refresh-token-hash", nil)
		d.mockSessionRepository.EXPECT().GetByRefreshTokenHash(gomock.Any(), gomock.Any()).Return(ses, nil)
		d.mockUUIDService.EXPECT().Generate().Return("new-refresh-token")
		d.mockKeyedHashService.EXPECT().Generate(gomock.Any(), gomock.Any()).Return("new-refresh-token-hash", nil)
		d.mockSessionRepository.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, _ *session.Session) error {
				trx.AddCommit(ctx, func() error {
					return errors.New("commit test error")
				})

				return nil
			})
		d.mockJWTService.EXPECT().Generate(gomock.Any(), gomock.Any()).Return("test-access-token", nil)

		res, err := d.usecase.Do(ctx, "old-refresh-token")
		assert.EqualError(t, err, "commit trx: commit test error")
		assert.Nil(t, res)
		assert.Empty(t, d.logbuf.Logs())
	})

	t.Run("generate input refresh token hash error", func(t *testing.T) {
		d := newDeps()

		d.mockTimeService.EXPECT().GetCurrentTime().Return(now)
		d.mockKeyedHashService.EXPECT().
			Generate(gomock.Eq([]byte("test-refresh-secret")), gomock.Eq("old-refresh-token")).
			Return("", errors.New("test error"))

		res, err := d.usecase.Do(ctx, "old-refresh-token")
		assert.EqualError(t, err, "generate input refresh token hash: test error")
		assert.Nil(t, res)
	})

	t.Run("invalid token", func(t *testing.T) {
		d := newDeps()

		d.mockTimeService.EXPECT().GetCurrentTime().Return(now)
		d.mockKeyedHashService.EXPECT().Generate(gomock.Any(), gomock.Any()).Return("old-refresh-token-hash", nil)
		d.mockSessionRepository.EXPECT().
			GetByRefreshTokenHash(gomock.Any(), gomock.Eq("old-refresh-token-hash")).
			Return(nil, repository.ErrNotFound)

		res, err := d.usecase.Do(ctx, "old-refresh-token")
		assert.EqualError(t, err, "invalid token")
		assert.ErrorIs(t, err, dto.ErrInvalidToken)
		assert.Nil(t, res)
		assert.Empty(t, d.logbuf.Logs())
	})

	t.Run("get by refresh token error + trx rollback error", func(t *testing.T) {
		d := newDeps()

		d.mockTimeService.EXPECT().GetCurrentTime().Return(now)
		d.mockKeyedHashService.EXPECT().Generate(gomock.Any(), gomock.Any()).Return("old-refresh-token-hash", nil)
		d.mockSessionRepository.EXPECT().
			GetByRefreshTokenHash(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, _ string) (*session.Session, error) {
				trx.AddRollback(ctx, func() error {
					return errors.New("rollback test error")
				})

				return nil, errors.New("test error")
			})

		res, err := d.usecase.Do(ctx, "old-refresh-token")
		assert.EqualError(t, err, "get by refresh token: test error")
		assert.Nil(t, res)

		logs := d.logbuf.Logs()
		assert.Len(t, logs, 1)
		assert.JSONEq(t, `{
			"created_at": "test-created-at",
			"level": "error",
			"pkg": "internal/app/usecases/refresh_session",
			"message": "trx rollback error",
			"error": "rollback test error",
			"additional_info": "get by refresh token: test error"
		}`, logs[0])
	})

	t.Run("generate refresh token hash error + trx rollback error", func(t *testing.T) {
		d := newDeps()

		d.mockTimeService.EXPECT().GetCurrentTime().Return(now)
		d.mockKeyedHashService.EXPECT().Generate(gomock.Any(), gomock.Any()).Return("old-refresh-token-hash", nil)
		d.mockSessionRepository.EXPECT().
			GetByRefreshTokenHash(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, _ string) (*session.Session, error) {
				trx.AddRollback(ctx, func() error {
					return errors.New("rollback test error")
				})

				return &session.Session{ID: "test-session-id", UserID: "test-user-id"}, nil
			})
		d.mockUUIDService.EXPECT().Generate().Return("new-refresh-token")
		d.mockKeyedHashService.EXPECT().Generate(gomock.Any(), gomock.Eq("new-refresh-token")).Return("", errors.New("test error"))

		res, err := d.usecase.Do(ctx, "old-refresh-token")
		assert.EqualError(t, err, "generate refresh token hash: test error")
		assert.Nil(t, res)

		logs := d.logbuf.Logs()
		assert.Len(t, logs, 1)
		assert.JSONEq(t, `{
			"created_at": "test-created-at",
			"level": "error",
			"pkg": "internal/app/usecases/refresh_session",
			"message": "trx rollback error",
			"error": "rollback test error",
			"additional_info": "generate refresh token hash: test error"
		}`, logs[0])
	})

	t.Run("save session error + trx rollback error", func(t *testing.T) {
		d := newDeps()

		d.mockTimeService.EXPECT().GetCurrentTime().Return(now)
		d.mockKeyedHashService.EXPECT().Generate(gomock.Any(), gomock.Any()).Return("old-refresh-token-hash", nil)
		d.mockSessionRepository.EXPECT().
			GetByRefreshTokenHash(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, _ string) (*session.Session, error) {
				trx.AddRollback(ctx, func() error {
					return errors.New("rollback test error")
				})

				return &session.Session{ID: "test-session-id", UserID: "test-user-id"}, nil
			})
		d.mockUUIDService.EXPECT().Generate().Return("new-refresh-token")
		d.mockKeyedHashService.EXPECT().Generate(gomock.Any(), gomock.Eq("new-refresh-token")).Return("new-refresh-token-hash", nil)
		d.mockSessionRepository.EXPECT().Save(gomock.Any(), gomock.Any()).Return(errors.New("test error"))

		res, err := d.usecase.Do(ctx, "old-refresh-token")
		assert.EqualError(t, err, "save session: test error")
		assert.Nil(t, res)

		logs := d.logbuf.Logs()
		assert.Len(t, logs, 1)
		assert.JSONEq(t, `{
			"created_at": "test-created-at",
			"level": "error",
			"pkg": "internal/app/usecases/refresh_session",
			"message": "trx rollback error",
			"error": "rollback test error",
			"additional_info": "save session: test error"
		}`, logs[0])
	})

	t.Run("generate access token error + trx rollback error", func(t *testing.T) {
		d := newDeps()

		d.mockTimeService.EXPECT().GetCurrentTime().Return(now)
		d.mockKeyedHashService.EXPECT().Generate(gomock.Any(), gomock.Any()).Return("old-refresh-token-hash", nil)
		d.mockSessionRepository.EXPECT().
			GetByRefreshTokenHash(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, _ string) (*session.Session, error) {
				trx.AddRollback(ctx, func() error {
					return errors.New("rollback test error")
				})

				return &session.Session{ID: "test-session-id", UserID: "test-user-id"}, nil
			})
		d.mockUUIDService.EXPECT().Generate().Return("new-refresh-token")
		d.mockKeyedHashService.EXPECT().Generate(gomock.Any(), gomock.Eq("new-refresh-token")).Return("new-refresh-token-hash", nil)
		d.mockSessionRepository.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
		d.mockJWTService.EXPECT().Generate(gomock.Any(), gomock.Any()).Return("", errors.New("test error"))

		res, err := d.usecase.Do(ctx, "old-refresh-token")
		assert.EqualError(t, err, "generate access token as jwt: test error")
		assert.Nil(t, res)

		logs := d.logbuf.Logs()
		assert.Len(t, logs, 1)
		assert.JSONEq(t, `{
			"created_at": "test-created-at",
			"level": "error",
			"pkg": "internal/app/usecases/refresh_session",
			"message": "trx rollback error",
			"error": "rollback test error",
			"additional_info": "generate access token as jwt: test error"
		}`, logs[0])
	})
}
