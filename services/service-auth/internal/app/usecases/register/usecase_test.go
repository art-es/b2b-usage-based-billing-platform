package register

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/repository"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/register/dto"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/trx"
)

func TestUsecase(t *testing.T) {
	ctx := context.Background()

	type deps struct {
		mockUserRepository              *MockuserRepository
		mockHashService                 *MockhashService
		mockEmailVerificationRepository *MockemailVerificationRepository
		logbuf                          log.Buffer
		usecase                         *Usecase
	}

	newDeps := func() *deps {
		mockCtrl := gomock.NewController(t)
		mockHashService := NewMockhashService(mockCtrl)
		mockUserRepository := NewMockuserRepository(mockCtrl)
		mockVerificationRepository := NewMockemailVerificationRepository(mockCtrl)

		logbuf := log.NewBuffer()
		logger := log.NewLogger(&log.Options{
			Output:       logbuf,
			GetCreatedAt: func() string { return "test-created-at" },
		})

		return &deps{
			mockUserRepository:              mockUserRepository,
			mockHashService:                 mockHashService,
			mockEmailVerificationRepository: mockVerificationRepository,
			logbuf:                          logbuf,
			usecase: NewUsecase(
				mockHashService,
				mockUserRepository,
				mockVerificationRepository,
				logger,
			),
		}
	}

	t.Run("ok", func(t *testing.T) {
		d := newDeps()

		d.mockHashService.EXPECT().
			Generate(gomock.Eq("test-password")).
			Return("test-password-hash", nil)

		expUser := &user.User{
			Name:         "test-name",
			Email:        "test-email",
			PasswordHash: "test-password-hash",
		}

		d.mockUserRepository.EXPECT().
			Create(gomock.Any(), gomock.Eq(expUser)).
			Do(func(_ context.Context, u *user.User) {
				u.ID = "test-user-id"
			}).
			Return(nil)

		d.mockEmailVerificationRepository.EXPECT().
			Create(gomock.Any(), gomock.Eq("test-user-id")).
			Return(nil)

		req := &dto.Request{
			Name:     "test-name",
			Email:    "test-email",
			Password: "test-password",
		}

		err := d.usecase.Do(ctx, req)
		assert.NoError(t, err)
		assert.Empty(t, d.logbuf.Logs())
	})

	t.Run("commit trx", func(t *testing.T) {
		d := newDeps()

		d.mockHashService.EXPECT().
			Generate(gomock.Eq("test-password")).
			Return("test-password-hash", nil)

		expUser := &user.User{
			Name:         "test-name",
			Email:        "test-email",
			PasswordHash: "test-password-hash",
		}

		d.mockUserRepository.EXPECT().
			Create(gomock.Any(), gomock.Eq(expUser)).
			Do(func(ctx context.Context, u *user.User) {
				u.ID = "test-user-id"

				trx.AddCommit(ctx, func() error {
					return errors.New("commit test error")
				})
			}).
			Return(nil)

		d.mockEmailVerificationRepository.EXPECT().
			Create(gomock.Any(), gomock.Eq("test-user-id")).
			Return(nil)

		req := &dto.Request{
			Name:     "test-name",
			Email:    "test-email",
			Password: "test-password",
		}

		err := d.usecase.Do(ctx, req)
		assert.EqualError(t, err, "commit trx: commit test error")
		assert.Empty(t, d.logbuf.Logs())
	})

	t.Run("create verification error + trx rollback error", func(t *testing.T) {
		d := newDeps()

		d.mockHashService.EXPECT().
			Generate(gomock.Eq("test-password")).
			Return("test-password-hash", nil)

		expUser := &user.User{
			Name:         "test-name",
			Email:        "test-email",
			PasswordHash: "test-password-hash",
		}

		d.mockUserRepository.EXPECT().
			Create(gomock.Any(), gomock.Eq(expUser)).
			Do(func(ctx context.Context, u *user.User) {
				u.ID = "test-user-id"

				trx.AddRollback(ctx, func() error {
					return errors.New("rollback test error")
				})
			}).
			Return(nil)

		d.mockEmailVerificationRepository.EXPECT().
			Create(gomock.Any(), gomock.Eq("test-user-id")).
			Return(errors.New("test error"))

		req := &dto.Request{
			Name:     "test-name",
			Email:    "test-email",
			Password: "test-password",
		}

		err := d.usecase.Do(ctx, req)
		assert.EqualError(t, err, "create verification: test error")

		logs := d.logbuf.Logs()
		assert.Len(t, logs, 1)
		assert.JSONEq(t, `{
			"created_at": "test-created-at",
			"level": "error",
			"pkg": "internal/app/usecases/register",
			"message": "trx rollback error",
			"error": "rollback test error",
			"additional_info": "create verification: test error"
		}`, logs[0])
	})

	t.Run("create user error", func(t *testing.T) {
		d := newDeps()

		d.mockHashService.EXPECT().
			Generate(gomock.Eq("test-password")).
			Return("test-password-hash", nil)

		expUser := &user.User{
			Name:         "test-name",
			Email:        "test-email",
			PasswordHash: "test-password-hash",
		}

		d.mockUserRepository.EXPECT().
			Create(gomock.Any(), gomock.Eq(expUser)).
			Return(errors.New("test error"))

		req := &dto.Request{
			Name:     "test-name",
			Email:    "test-email",
			Password: "test-password",
		}

		err := d.usecase.Do(ctx, req)
		assert.EqualError(t, err, "create user: test error")
		assert.Empty(t, d.logbuf.Logs())
	})

	t.Run("email is already in use", func(t *testing.T) {
		d := newDeps()

		d.mockHashService.EXPECT().
			Generate(gomock.Eq("test-password")).
			Return("test-password-hash", nil)

		expUser := &user.User{
			Name:         "test-name",
			Email:        "test-email",
			PasswordHash: "test-password-hash",
		}

		d.mockUserRepository.EXPECT().
			Create(gomock.Any(), gomock.Eq(expUser)).
			Return(fmt.Errorf("repo: %w", repository.ErrUnique))

		req := &dto.Request{
			Name:     "test-name",
			Email:    "test-email",
			Password: "test-password",
		}

		err := d.usecase.Do(ctx, req)
		assert.EqualError(t, err, "email is already in use")
		assert.ErrorIs(t, err, dto.ErrEmailInUse)
		assert.Empty(t, d.logbuf.Logs())
	})

	t.Run("generate password hash error", func(t *testing.T) {
		d := newDeps()

		d.mockHashService.EXPECT().
			Generate(gomock.Eq("test-password")).
			Return("", errors.New("test error"))

		req := &dto.Request{
			Name:     "test-name",
			Email:    "test-email",
			Password: "test-password",
		}

		err := d.usecase.Do(ctx, req)
		assert.EqualError(t, err, "generate password hash: test error")
		assert.Empty(t, d.logbuf.Logs())
	})
}
