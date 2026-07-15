package verify_email

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/verify_email/dto"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/trx"
)

func TestUsecase(t *testing.T) {
	ctx := context.Background()

	type deps struct {
		mockEmailVerificationsRepository *MockemailVerificationsRepository
		mockUserRepository               *MockuserRepository
		logbuf                           log.Buffer
		usecase                          *Usecase
	}

	newDeps := func() *deps {
		mockCtrl := gomock.NewController(t)
		mockEmailVerificationsRepository := NewMockemailVerificationsRepository(mockCtrl)
		mockUserRepository := NewMockuserRepository(mockCtrl)

		logbuf := log.NewBuffer()
		logger := log.NewLogger(&log.Options{
			Output:       logbuf,
			GetCreatedAt: func() string { return "test-created-at" },
		})

		return &deps{
			mockEmailVerificationsRepository: mockEmailVerificationsRepository,
			mockUserRepository:               mockUserRepository,
			logbuf:                           logbuf,
			usecase: NewUsecase(
				mockEmailVerificationsRepository,
				mockUserRepository,
				logger,
			),
		}
	}

	t.Run("ok", func(t *testing.T) {
		d := newDeps()

		d.mockEmailVerificationsRepository.EXPECT().
			GetByToken(gomock.Any(), gomock.Eq("test-token")).
			Return(&user.EmailVerification{UserID: "test-user-id"}, nil)

		d.mockUserRepository.EXPECT().
			MarkAsVerified(gomock.Any(), gomock.Eq("test-user-id")).
			Return(nil)

		d.mockEmailVerificationsRepository.EXPECT().
			DeleteTokensByUserID(gomock.Any(), gomock.Eq("test-user-id")).
			Return(nil)

		err := d.usecase.Do(ctx, "test-token")
		assert.NoError(t, err)
		assert.Empty(t, d.logbuf.Logs())
	})

	t.Run("commit trx", func(t *testing.T) {
		d := newDeps()

		d.mockEmailVerificationsRepository.EXPECT().
			GetByToken(gomock.Any(), gomock.Eq("test-token")).
			Return(&user.EmailVerification{UserID: "test-user-id"}, nil)

		d.mockUserRepository.EXPECT().
			MarkAsVerified(gomock.Any(), gomock.Eq("test-user-id")).
			Return(nil)

		d.mockEmailVerificationsRepository.EXPECT().
			DeleteTokensByUserID(gomock.Any(), gomock.Eq("test-user-id")).
			DoAndReturn(func(ctx context.Context, _ string) error {
				trx.AddCommit(ctx, func() error {
					return errors.New("commit test error")
				})

				return nil
			})

		err := d.usecase.Do(ctx, "test-token")
		assert.EqualError(t, err, "commit trx: commit test error")
		assert.Empty(t, d.logbuf.Logs())
	})

	t.Run("get verification error + trx rollback error", func(t *testing.T) {
		d := newDeps()

		d.mockEmailVerificationsRepository.EXPECT().
			GetByToken(gomock.Any(), gomock.Eq("test-token")).
			DoAndReturn(func(ctx context.Context, _ string) (*user.EmailVerification, error) {
				trx.AddRollback(ctx, func() error {
					return errors.New("rollback test error")
				})

				return nil, errors.New("test error")
			})

		err := d.usecase.Do(ctx, "test-token")
		assert.EqualError(t, err, "get verification by token: test error")

		logs := d.logbuf.Logs()
		assert.Len(t, logs, 1)
		assert.JSONEq(t, `{
			"created_at": "test-created-at",
			"level": "error",
			"message": "trx rollback error",
			"error": "rollback test error",
			"additional_info": "get verification by token: test error"
		}`, logs[0])
	})

	t.Run("invalid token", func(t *testing.T) {
		d := newDeps()

		d.mockEmailVerificationsRepository.EXPECT().
			GetByToken(gomock.Any(), gomock.Eq("test-token")).
			Return(nil, nil)

		err := d.usecase.Do(ctx, "test-token")
		assert.EqualError(t, err, "invalid token")
		assert.ErrorIs(t, err, dto.ErrInvalidToken)
		assert.Empty(t, d.logbuf.Logs())
	})

	t.Run("mark user as verified error", func(t *testing.T) {
		d := newDeps()

		d.mockEmailVerificationsRepository.EXPECT().
			GetByToken(gomock.Any(), gomock.Eq("test-token")).
			Return(&user.EmailVerification{UserID: "test-user-id"}, nil)

		d.mockUserRepository.EXPECT().
			MarkAsVerified(gomock.Any(), gomock.Eq("test-user-id")).
			Return(errors.New("test error"))

		err := d.usecase.Do(ctx, "test-token")
		assert.EqualError(t, err, "mark user as verified: test error")
		assert.Empty(t, d.logbuf.Logs())
	})

	t.Run("delete tokens error", func(t *testing.T) {
		d := newDeps()

		d.mockEmailVerificationsRepository.EXPECT().
			GetByToken(gomock.Any(), gomock.Eq("test-token")).
			Return(&user.EmailVerification{UserID: "test-user-id"}, nil)

		d.mockUserRepository.EXPECT().
			MarkAsVerified(gomock.Any(), gomock.Eq("test-user-id")).
			Return(nil)

		d.mockEmailVerificationsRepository.EXPECT().
			DeleteTokensByUserID(gomock.Any(), gomock.Eq("test-user-id")).
			Return(errors.New("test error"))

		err := d.usecase.Do(ctx, "test-token")
		assert.EqualError(t, err, "delete verification tokens by user id: test error")
		assert.Empty(t, d.logbuf.Logs())
	})
}
