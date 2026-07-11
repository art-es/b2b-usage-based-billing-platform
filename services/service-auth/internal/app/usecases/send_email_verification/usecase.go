//go:generate mockgen -source=usecase.go -destination=usecase_mock_test.go -package=$GOPACKAGE
package send_email_verification

import (
	"context"
	"errors"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/event"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/trx"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/trx/trxutil"
)

type verificationRepository interface {
	GetUnsent(ctx context.Context, batchSize int) ([]*user.EmailVerification, error)
	MarkAsSent(ctx context.Context, tokens []string) error
}

type emailSendProducer interface {
	Produce(ctx context.Context, ev event.EmailSend) error
}

type Usecase struct {
	verificationRepository verificationRepository
	emailSendProducer      emailSendProducer
	logger                 log.Logger
	batchSize              int
}

func NewUsecase(
	verificationRepository verificationRepository,
	emailSendProducer emailSendProducer,
	logger log.Logger,
	batchSize int,
) *Usecase {
	logger = logger.Set("pkg", "internal/app/usecases/mail_verification")

	return &Usecase{
		verificationRepository: verificationRepository,
		emailSendProducer:      emailSendProducer,
		logger:                 logger,
		batchSize:              batchSize,
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
	vers, err := u.verificationRepository.GetUnsent(ctx, u.batchSize)
	if err != nil {
		return 0, fmt.Errorf("get unsent verifications: %w", err)
	}

	if len(vers) == 0 {
		return 0, nil
	}

	vers = u.produceEmailSend(ctx, vers)

	if len(vers) == 0 {
		return 0, errors.New("all email.send producing failed")
	}

	err = u.verificationRepository.MarkAsSent(ctx, convertToTokens(vers))
	if err != nil {
		return 0, fmt.Errorf("mark verifications as sent: %w", err)
	}

	return len(vers), nil
}

func (u *Usecase) produceEmailSend(ctx context.Context, vers []*user.EmailVerification) []*user.EmailVerification {
	sent := make([]*user.EmailVerification, 0, len(vers))
	var unsentErr error

	for _, ver := range vers {
		err := u.emailSendProducer.Produce(ctx, convertToEvent(ver))

		if err != nil {
			unsentErr = errors.Join(unsentErr, err)
		} else {
			sent = append(sent, ver)
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

func convertToEvent(ver *user.EmailVerification) event.EmailSend {
	return event.EmailSend{
		IdempotencyKey: "email-verification:" + ver.Token,
		Email:          ver.Email,
		Subject:        ver.EmailSubject(),
		Content:        ver.EmailContent(),
	}
}

func convertToTokens(vers []*user.EmailVerification) []string {
	tokens := make([]string, 0, len(vers))
	for _, ver := range vers {
		tokens = append(tokens, ver.Token)
	}
	return tokens
}
