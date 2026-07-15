package login

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/hash"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/session"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/repository"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/login/dto"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/trx"
)

func TestUsecase(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 7, 15, 12, 0, 0, 0, time.UTC)

	type deps struct {
		mockJWTService          *MockjwtService
		mockKeyedHashService    *MockkeyedHashService
		mockPasswordHashService *MockpasswordHashService
		mockTimeService         *MocktimeService
		mockUUIDService         *MockuuidService
		mockSessionRepository   *MocksessionRepository
		mockUserRepository      *MockuserRepository
		logbuf                  log.Buffer
		usecase                 *Usecase
	}

	newDeps := func() *deps {
		mockCtrl := gomock.NewController(t)
		mockJWTService := NewMockjwtService(mockCtrl)
		mockKeyedHashService := NewMockkeyedHashService(mockCtrl)
		mockPasswordHashService := NewMockpasswordHashService(mockCtrl)
		mockTimeService := NewMocktimeService(mockCtrl)
		mockUUIDService := NewMockuuidService(mockCtrl)
		mockSessionRepository := NewMocksessionRepository(mockCtrl)
		mockUserRepository := NewMockuserRepository(mockCtrl)

		logbuf := log.NewBuffer()
		logger := log.NewLogger(&log.Options{
			Output:       logbuf,
			GetCreatedAt: func() string { return "test-created-at" },
		})

		return &deps{
			mockJWTService:          mockJWTService,
			mockKeyedHashService:    mockKeyedHashService,
			mockPasswordHashService: mockPasswordHashService,
			mockTimeService:         mockTimeService,
			mockUUIDService:         mockUUIDService,
			mockSessionRepository:   mockSessionRepository,
			mockUserRepository:      mockUserRepository,
			logbuf:                  logbuf,
			usecase: NewUsecase(
				mockJWTService,
				mockKeyedHashService,
				mockPasswordHashService,
				mockTimeService,
				mockUUIDService,
				mockSessionRepository,
				mockUserRepository,
				"test-jwt-secret",
				"test-refresh-secret",
				logger,
			),
		}
	}

	t.Run("ok", func(t *testing.T) {
		d := newDeps()

		usr := &user.User{
			ID:           "test-user-id",
			Email:        "test@example.com",
			PasswordHash: "test-password-hash",
			IsVerified:   true,
		}

		d.mockTimeService.EXPECT().
			GetCurrentTime().
			Return(now)

		d.mockUserRepository.EXPECT().
			FindByEmail(gomock.Any(), gomock.Eq("test@example.com")).
			Return(usr, nil)

		d.mockPasswordHashService.EXPECT().
			Compare(gomock.Eq("test-password"), gomock.Eq("test-password-hash")).
			Return(nil)

		d.mockUUIDService.EXPECT().
			Generate().
			Return("test-refresh-token")

		d.mockKeyedHashService.EXPECT().
			Generate(gomock.Eq([]byte("test-refresh-secret")), gomock.Eq("test-refresh-token")).
			Return("test-refresh-token-hash", nil)

		expSession := session.NewSession("test-user-id", "test-refresh-token-hash", now)
		d.mockSessionRepository.EXPECT().
			Save(gomock.Any(), gomock.Eq(expSession)).
			Do(func(_ context.Context, ses *session.Session) {
				ses.ID = "test-session-id"
			}).
			Return(nil)

		d.mockJWTService.EXPECT().
			Generate(gomock.Eq([]byte("test-jwt-secret")), gomock.Eq(jwt.NewClaims("test-session-id", "test-user-id", nil))).
			Return("test-access-token", nil)

		res, err := d.usecase.Do(ctx, &dto.Request{
			Email:    "test@example.com",
			Password: "test-password",
		})
		assert.NoError(t, err)
		assert.Equal(t, &dto.Response{
			AccessToken:  "test-access-token",
			RefreshToken: "test-refresh-token",
		}, res)
		assert.Empty(t, d.logbuf.Logs())
	})

	t.Run("commit trx", func(t *testing.T) {
		d := newDeps()

		usr := &user.User{
			ID:           "test-user-id",
			Email:        "test@example.com",
			PasswordHash: "test-password-hash",
			IsVerified:   true,
		}

		d.mockTimeService.EXPECT().GetCurrentTime().Return(now)
		d.mockUserRepository.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(usr, nil)
		d.mockPasswordHashService.EXPECT().Compare(gomock.Any(), gomock.Any()).Return(nil)
		d.mockUUIDService.EXPECT().Generate().Return("test-refresh-token")
		d.mockKeyedHashService.EXPECT().Generate(gomock.Any(), gomock.Any()).Return("test-refresh-token-hash", nil)
		d.mockSessionRepository.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			Do(func(ctx context.Context, ses *session.Session) {
				ses.ID = "test-session-id"
				trx.AddCommit(ctx, func() error {
					return errors.New("commit test error")
				})
			}).
			Return(nil)
		d.mockJWTService.EXPECT().Generate(gomock.Any(), gomock.Any()).Return("test-access-token", nil)

		res, err := d.usecase.Do(ctx, &dto.Request{})
		assert.EqualError(t, err, "commit trx: commit test error")
		assert.Nil(t, res)
		assert.Empty(t, d.logbuf.Logs())
	})

	t.Run("user not found", func(t *testing.T) {
		d := newDeps()

		d.mockTimeService.EXPECT().GetCurrentTime().Return(now)
		d.mockUserRepository.EXPECT().
			FindByEmail(gomock.Any(), gomock.Eq("test@example.com")).
			Return(nil, repository.ErrNotFound)

		res, err := d.usecase.Do(ctx, &dto.Request{Email: "test@example.com"})
		assert.EqualError(t, err, "wrong credentials")
		assert.ErrorIs(t, err, dto.ErrWrongCredentials)
		assert.Nil(t, res)
	})

	t.Run("find user error", func(t *testing.T) {
		d := newDeps()

		d.mockTimeService.EXPECT().GetCurrentTime().Return(now)
		d.mockUserRepository.EXPECT().
			FindByEmail(gomock.Any(), gomock.Eq("test@example.com")).
			Return(nil, errors.New("test error"))

		res, err := d.usecase.Do(ctx, &dto.Request{Email: "test@example.com"})
		assert.EqualError(t, err, "get user by email: test error")
		assert.Nil(t, res)
	})

	t.Run("password mismatch", func(t *testing.T) {
		d := newDeps()

		usr := &user.User{PasswordHash: "test-password-hash", IsVerified: true}

		d.mockTimeService.EXPECT().GetCurrentTime().Return(now)
		d.mockUserRepository.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(usr, nil)
		d.mockPasswordHashService.EXPECT().
			Compare(gomock.Eq("test-password"), gomock.Eq("test-password-hash")).
			Return(hash.ErrMismatch)

		res, err := d.usecase.Do(ctx, &dto.Request{Password: "test-password"})
		assert.EqualError(t, err, "wrong credentials")
		assert.ErrorIs(t, err, dto.ErrWrongCredentials)
		assert.Nil(t, res)
	})

	t.Run("compare password error", func(t *testing.T) {
		d := newDeps()

		usr := &user.User{PasswordHash: "test-password-hash", IsVerified: true}

		d.mockTimeService.EXPECT().GetCurrentTime().Return(now)
		d.mockUserRepository.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(usr, nil)
		d.mockPasswordHashService.EXPECT().Compare(gomock.Any(), gomock.Any()).Return(errors.New("test error"))

		res, err := d.usecase.Do(ctx, &dto.Request{})
		assert.EqualError(t, err, "compare password with hash: test error")
		assert.Nil(t, res)
	})

	t.Run("email is not verified", func(t *testing.T) {
		d := newDeps()

		usr := &user.User{PasswordHash: "test-password-hash"}

		d.mockTimeService.EXPECT().GetCurrentTime().Return(now)
		d.mockUserRepository.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(usr, nil)
		d.mockPasswordHashService.EXPECT().Compare(gomock.Any(), gomock.Any()).Return(nil)

		res, err := d.usecase.Do(ctx, &dto.Request{})
		assert.EqualError(t, err, "email is not verified")
		assert.ErrorIs(t, err, dto.ErrEmailNotVerified)
		assert.Nil(t, res)
	})

	t.Run("generate refresh token hash error", func(t *testing.T) {
		d := newDeps()

		usr := &user.User{ID: "test-user-id", PasswordHash: "test-password-hash", IsVerified: true}

		d.mockTimeService.EXPECT().GetCurrentTime().Return(now)
		d.mockUserRepository.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(usr, nil)
		d.mockPasswordHashService.EXPECT().Compare(gomock.Any(), gomock.Any()).Return(nil)
		d.mockUUIDService.EXPECT().Generate().Return("test-refresh-token")
		d.mockKeyedHashService.EXPECT().
			Generate(gomock.Eq([]byte("test-refresh-secret")), gomock.Eq("test-refresh-token")).
			Return("", errors.New("test error"))

		res, err := d.usecase.Do(ctx, &dto.Request{})
		assert.EqualError(t, err, "generate refresh token hash: test error")
		assert.Nil(t, res)
	})

	t.Run("save session error + trx rollback error", func(t *testing.T) {
		d := newDeps()

		usr := &user.User{ID: "test-user-id", PasswordHash: "test-password-hash", IsVerified: true}

		d.mockTimeService.EXPECT().GetCurrentTime().Return(now)
		d.mockUserRepository.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(usr, nil)
		d.mockPasswordHashService.EXPECT().Compare(gomock.Any(), gomock.Any()).Return(nil)
		d.mockUUIDService.EXPECT().Generate().Return("test-refresh-token")
		d.mockKeyedHashService.EXPECT().Generate(gomock.Any(), gomock.Any()).Return("test-refresh-token-hash", nil)
		d.mockSessionRepository.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, _ *session.Session) error {
				trx.AddRollback(ctx, func() error {
					return errors.New("rollback test error")
				})

				return errors.New("test error")
			})

		res, err := d.usecase.Do(ctx, &dto.Request{})
		assert.EqualError(t, err, "save session: test error")
		assert.Nil(t, res)

		logs := d.logbuf.Logs()
		assert.Len(t, logs, 1)
		assert.JSONEq(t, `{
			"created_at": "test-created-at",
			"level": "error",
			"pkg": "internal/app/usecases/login",
			"message": "trx rollback error",
			"error": "rollback test error",
			"additional_info": "save session: test error"
		}`, logs[0])
	})

	t.Run("generate access token error + trx rollback error", func(t *testing.T) {
		d := newDeps()

		usr := &user.User{ID: "test-user-id", PasswordHash: "test-password-hash", IsVerified: true}

		d.mockTimeService.EXPECT().GetCurrentTime().Return(now)
		d.mockUserRepository.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(usr, nil)
		d.mockPasswordHashService.EXPECT().Compare(gomock.Any(), gomock.Any()).Return(nil)
		d.mockUUIDService.EXPECT().Generate().Return("test-refresh-token")
		d.mockKeyedHashService.EXPECT().Generate(gomock.Any(), gomock.Any()).Return("test-refresh-token-hash", nil)
		d.mockSessionRepository.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			Do(func(ctx context.Context, ses *session.Session) {
				ses.ID = "test-session-id"
				trx.AddRollback(ctx, func() error {
					return errors.New("rollback test error")
				})
			}).
			Return(nil)
		d.mockJWTService.EXPECT().Generate(gomock.Any(), gomock.Any()).Return("", errors.New("test error"))

		res, err := d.usecase.Do(ctx, &dto.Request{})
		assert.EqualError(t, err, "generate access token as jwt: test error")
		assert.Nil(t, res)

		logs := d.logbuf.Logs()
		assert.Len(t, logs, 1)
		assert.JSONEq(t, `{
			"created_at": "test-created-at",
			"level": "error",
			"pkg": "internal/app/usecases/login",
			"message": "trx rollback error",
			"error": "rollback test error",
			"additional_info": "generate access token as jwt: test error"
		}`, logs[0])
	})
}
