package send_password_reset

import (
	"context"
	"errors"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/event"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/trx"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/trx/trxutil"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
)

type passwordResetRepository interface {
	GetUnsent(ctx context.Context, batchSize int) ([]*user.PasswordReset, error)
	MarkAsSent(ctx context.Context, tokens []string) error
}

type emailSendProducer interface {
	Produce(ctx context.Context, ev event.EmailSend) error
}

type Usecase struct {
	passwordResetRepository passwordResetRepository
	emailSendProducer       emailSendProducer
	logger                  log.Logger
	batchSize               int
}

func NewUsecase(
	passwordResetRepository passwordResetRepository,
	emailSendProducer emailSendProducer,
	logger log.Logger,
	batchSize int,
) *Usecase {
	logger = logger.Set("pkg", "internal/app/usecases/send_password_reset")

	return &Usecase{
		passwordResetRepository: passwordResetRepository,
		emailSendProducer:       emailSendProducer,
		logger:                  logger,
		batchSize:               batchSize,
	}
}

func (u *Usecase) Do(ctx context.Context) (int, error) {
	ctx = trx.Begin(ctx)

	updatesCount, err := u.processTrx(ctx)
	if err != nil {
		trxutil.RollbackOrLog(ctx, u.logger, err.Error())

		return 0, err
	}

	err = trx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("commit trx: %w", err)
	}

	return updatesCount, nil
}

func (u *Usecase) processTrx(ctx context.Context) (int, error) {
	resets, err := u.passwordResetRepository.GetUnsent(ctx, u.batchSize)
	if err != nil {
		return 0, fmt.Errorf("get unsent resets: %w", err)
	}

	if len(resets) == 0 {
		return 0, nil
	}

	resets = u.produceEmailSend(ctx, resets)

	if len(resets) == 0 {
		return 0, errors.New("all email.send producing failed")
	}

	err = u.passwordResetRepository.MarkAsSent(ctx, convertToTokens(resets))
	if err != nil {
		return 0, fmt.Errorf("mark resets as sent: %w", err)
	}

	return len(resets), nil
}

func (u *Usecase) produceEmailSend(ctx context.Context, resets []*user.PasswordReset) []*user.PasswordReset {
	sent := make([]*user.PasswordReset, 0, len(resets))
	var unsentErr error

	for _, reset := range resets {
		err := u.emailSendProducer.Produce(ctx, convertToEvent(reset))

		if err != nil {
			unsentErr = errors.Join(unsentErr, err)
		} else {
			sent = append(sent, reset)
		}
	}

	if unsentErr != nil {
		u.logger.Log(log.Error).
			Set("message", "produce email.send error").
			Set("error", unsentErr.Error()).
			Write()
	}

	return sent
}

func convertToEvent(reset *user.PasswordReset) event.EmailSend {
	return event.EmailSend{
		IdempotencyKey: "password-reset:" + reset.Token,
		Email:          reset.Email,
		Subject:        reset.EmailSubject(),
		Content:        reset.EmailContent(),
	}
}

func convertToTokens(resets []*user.PasswordReset) []string {
	tokens := make([]string, 0, len(resets))
	for _, reset := range resets {
		tokens = append(tokens, reset.Token)
	}
	return tokens
}
