package email_verification

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
	GetUnsent(ctx context.Context, batchSize int) ([]*user.Verification, error)
	MarkAsSent(ctx context.Context, tokens []string) error
}

type emailSendProducer interface {
	Produce(ctx context.Context, events []event.EmailSend) error
}

type Usecase struct {
	verificationRepository verificationRepository
	emailEventProducer     emailSendProducer
	logger                 log.Logger
	batchSize              int
}

func NewUsecase(
	verificationRepository verificationRepository,
	emailSendEventProducer emailSendProducer,
	logger log.Logger,
	batchSize int,
) *Usecase {
	logger = logger.Set("pkg", "internal/app/usecases/mail_verification")

	return &Usecase{
		verificationRepository: verificationRepository,
		emailEventProducer:     emailSendEventProducer,
		logger:                 logger,
		batchSize:              batchSize,
	}
}

func (u *Usecase) Do(ctx context.Context) (int, error) {
	ctx = trx.Begin(ctx)

	vers, err := u.verificationRepository.GetUnsent(ctx, u.batchSize)
	if err != nil {
		trxutil.RollbackOrLog(ctx, u.logger, fmt.Sprintf("get unsent verifications: %v", err))

		return 0, fmt.Errorf("get unsent verifications: %w", err)
	}

	if len(vers) == 0 {
		return 0, nil
	}

	err = u.emailEventProducer.Produce(ctx, convertToEvents(vers))
	if err != nil {
		trxutil.RollbackOrLog(ctx, u.logger, fmt.Sprintf("publish mails: %v", err))

		return 0, fmt.Errorf("publish mails: %w", err)
	}

	err = u.verificationRepository.MarkAsSent(ctx, convertToTokens(vers))
	if err != nil {
		trxutil.RollbackOrLog(ctx, u.logger, fmt.Sprintf("mark verifications as sent: %v", err))

		return 0, fmt.Errorf("mark verifications as sent: %w", err)
	}

	err = trx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("commit trx: %w", err)
	}

	return len(vers), nil
}

func convertToEvents(vers []*user.Verification) []event.EmailSend {
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

func convertToTokens(vers []*user.Verification) []string {
	tokens := make([]string, 0, len(vers))
	for _, ver := range vers {
		tokens = append(tokens, ver.Token)
	}
	return tokens
}
