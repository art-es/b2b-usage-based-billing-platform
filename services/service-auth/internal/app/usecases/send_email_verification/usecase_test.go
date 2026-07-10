package send_email_verification

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/event"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/trx"
)

func TestUsecase(t *testing.T) {
	ctx := context.Background()

	type deps struct {
		mockVerificationRepository *MockverificationRepository
		mockEmailSendProducer      *MockemailSendProducer
		logbuf                     log.Buffer
		usecase                    *Usecase
	}

	newDeps := func() *deps {
		mockCtrl := gomock.NewController(t)
		mockVerificationRepository := NewMockverificationRepository(mockCtrl)
		mockEmailSendProducer := NewMockemailSendProducer(mockCtrl)

		logbuf := log.NewBuffer()
		logger := log.NewLogger(&log.Options{
			Output:       logbuf,
			GetCreatedAt: func() string { return "test-created-at" },
		})

		return &deps{
			mockVerificationRepository: mockVerificationRepository,
			mockEmailSendProducer:      mockEmailSendProducer,
			logbuf:                     logbuf,
			usecase: NewUsecase(
				mockVerificationRepository,
				mockEmailSendProducer,
				logger,
				10,
			),
		}
	}

	t.Run("ok", func(t *testing.T) {
		d := newDeps()

		vers := []*user.EmailVerification{
			{
				Token: "test-token-1",
				Email: "test-1@example.com",
			},
			{
				Token: "test-token-2",
				Email: "test-2@example.com",
			},
		}

		d.mockVerificationRepository.EXPECT().
			GetUnsent(gomock.Any(), gomock.Eq(10)).
			Return(vers, nil)

		d.mockEmailSendProducer.EXPECT().
			Produce(gomock.Any(), gomock.Eq([]event.EmailSend{
				{
					Email:   "test-1@example.com",
					Subject: "Verify your new account",
					Content: "To verify your email address, please follow the link: https://example.com/verify-email/test-token-1",
				},
				{
					Email:   "test-2@example.com",
					Subject: "Verify your new account",
					Content: "To verify your email address, please follow the link: https://example.com/verify-email/test-token-2",
				},
			})).
			Return(nil)

		d.mockVerificationRepository.EXPECT().
			MarkAsSent(gomock.Any(), gomock.Eq([]string{"test-token-1", "test-token-2"})).
			Return(nil)

		updates, err := d.usecase.Do(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 2, updates)
		assert.Empty(t, d.logbuf.Logs())
	})

	t.Run("no unsent verifications", func(t *testing.T) {
		d := newDeps()

		commitCalled := false

		d.mockVerificationRepository.EXPECT().
			GetUnsent(gomock.Any(), gomock.Eq(10)).
			DoAndReturn(func(ctx context.Context, _ int) ([]*user.EmailVerification, error) {
				trx.AddCommit(ctx, func() error {
					commitCalled = true
					return nil
				})

				return nil, nil
			})

		updates, err := d.usecase.Do(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 0, updates)
		assert.True(t, commitCalled)
		assert.Empty(t, d.logbuf.Logs())
	})

	t.Run("commit trx", func(t *testing.T) {
		d := newDeps()

		vers := []*user.EmailVerification{
			{
				Token: "test-token",
				Email: "test@example.com",
			},
		}

		d.mockVerificationRepository.EXPECT().
			GetUnsent(gomock.Any(), gomock.Eq(10)).
			Return(vers, nil)

		d.mockEmailSendProducer.EXPECT().
			Produce(gomock.Any(), gomock.Any()).
			Return(nil)

		d.mockVerificationRepository.EXPECT().
			MarkAsSent(gomock.Any(), gomock.Eq([]string{"test-token"})).
			DoAndReturn(func(ctx context.Context, _ []string) error {
				trx.AddCommit(ctx, func() error {
					return errors.New("commit test error")
				})

				return nil
			})

		updates, err := d.usecase.Do(ctx)
		assert.EqualError(t, err, "commit trx: commit test error")
		assert.Equal(t, 0, updates)
		assert.Empty(t, d.logbuf.Logs())
	})

	t.Run("get unsent error + trx rollback error", func(t *testing.T) {
		d := newDeps()

		d.mockVerificationRepository.EXPECT().
			GetUnsent(gomock.Any(), gomock.Eq(10)).
			DoAndReturn(func(ctx context.Context, _ int) ([]*user.EmailVerification, error) {
				trx.AddRollback(ctx, func() error {
					return errors.New("rollback test error")
				})

				return nil, errors.New("test error")
			})

		updates, err := d.usecase.Do(ctx)
		assert.EqualError(t, err, "get unsent verifications: test error")
		assert.Equal(t, 0, updates)

		logs := d.logbuf.Logs()
		assert.Len(t, logs, 1)
		assert.JSONEq(t, `{
			"created_at": "test-created-at",
			"level": "error",
			"pkg": "internal/app/usecases/mail_verification",
			"message": "trx rollback error",
			"error": "rollback test error",
			"additional_info": "get unsent verifications: test error"
		}`, logs[0])
	})

	t.Run("publish mails error + trx rollback error", func(t *testing.T) {
		d := newDeps()

		vers := []*user.EmailVerification{
			{
				Token: "test-token",
				Email: "test@example.com",
			},
		}

		d.mockVerificationRepository.EXPECT().
			GetUnsent(gomock.Any(), gomock.Eq(10)).
			DoAndReturn(func(ctx context.Context, _ int) ([]*user.EmailVerification, error) {
				trx.AddRollback(ctx, func() error {
					return errors.New("rollback test error")
				})

				return vers, nil
			})

		d.mockEmailSendProducer.EXPECT().
			Produce(gomock.Any(), gomock.Any()).
			Return(errors.New("test error"))

		updates, err := d.usecase.Do(ctx)
		assert.EqualError(t, err, "publish mails: test error")
		assert.Equal(t, 0, updates)

		logs := d.logbuf.Logs()
		assert.Len(t, logs, 1)
		assert.JSONEq(t, `{
			"created_at": "test-created-at",
			"level": "error",
			"pkg": "internal/app/usecases/mail_verification",
			"message": "trx rollback error",
			"error": "rollback test error",
			"additional_info": "publish mails: test error"
		}`, logs[0])
	})

	t.Run("mark as sent error + trx rollback error", func(t *testing.T) {
		d := newDeps()

		vers := []*user.EmailVerification{
			{
				Token: "test-token",
				Email: "test@example.com",
			},
		}

		d.mockVerificationRepository.EXPECT().
			GetUnsent(gomock.Any(), gomock.Eq(10)).
			DoAndReturn(func(ctx context.Context, _ int) ([]*user.EmailVerification, error) {
				trx.AddRollback(ctx, func() error {
					return errors.New("rollback test error")
				})

				return vers, nil
			})

		d.mockEmailSendProducer.EXPECT().
			Produce(gomock.Any(), gomock.Any()).
			Return(nil)

		d.mockVerificationRepository.EXPECT().
			MarkAsSent(gomock.Any(), gomock.Eq([]string{"test-token"})).
			Return(errors.New("test error"))

		updates, err := d.usecase.Do(ctx)
		assert.EqualError(t, err, "mark verifications as sent: test error")
		assert.Equal(t, 0, updates)

		logs := d.logbuf.Logs()
		assert.Len(t, logs, 1)
		assert.JSONEq(t, `{
			"created_at": "test-created-at",
			"level": "error",
			"pkg": "internal/app/usecases/mail_verification",
			"message": "trx rollback error",
			"error": "rollback test error",
			"additional_info": "mark verifications as sent: test error"
		}`, logs[0])
	})
}
