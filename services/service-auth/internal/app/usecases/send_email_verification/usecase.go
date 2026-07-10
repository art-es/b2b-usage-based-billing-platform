//go:generate mockgen -source=usecase.go -destination=usecase_mock_test.go -package=$GOPACKAGE
package send_email_verification

import (
	"context"
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
	Produce(ctx context.Context, events []event.EmailSend) error
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

	err = u.emailSendProducer.Produce(ctx, convertToEvents(vers))
	if err != nil {
		return 0, fmt.Errorf("produce email.send: %w", err)
	}

	err = u.verificationRepository.MarkAsSent(ctx, convertToTokens(vers))
	if err != nil {
		return 0, fmt.Errorf("mark verifications as sent: %w", err)
	}

	return len(vers), nil
}

func convertToEvents(vers []*user.EmailVerification) []event.EmailSend {
	events := make([]event.EmailSend, 0, len(vers))
	for _, ver := range vers {
		events = append(events, event.EmailSend{
			Email:   ver.Email,
			Subject: ver.EmailSubject(),
			Content: ver.EmailContent(),
		})
	}
	return events
}

func convertToTokens(vers []*user.EmailVerification) []string {
	tokens := make([]string, 0, len(vers))
	for _, ver := range vers {
		tokens = append(tokens, ver.Token)
	}
	return tokens
}
