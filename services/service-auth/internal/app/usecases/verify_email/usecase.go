package verify_email

import (
	"context"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/verify_email/dto"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/trx"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/trx/trxutil"
)

type emailVerificationsRepository interface {
	GetByToken(ctx context.Context, token string) (*user.EmailVerification, error)
	DeleteTokensByUserID(ctx context.Context, userID string) error
}

type userRepository interface {
	MarkAsVerified(ctx context.Context, userID string) error
}

type Usecase struct {
	emailVerificationsRepository emailVerificationsRepository
	userRepository               userRepository
	logger                       log.Logger
}

func NewUsecase(
	emailVerificationsRepository emailVerificationsRepository,
	userRepository userRepository,
	logger log.Logger,
) *Usecase {
	return &Usecase{
		emailVerificationsRepository: emailVerificationsRepository,
		userRepository:               userRepository,
		logger:                       logger,
	}
}

func (u *Usecase) Do(ctx context.Context, token string) error {
	ctx = trx.Begin(ctx)

	err := u.processTrx(ctx, token)
	if err != nil {
		trxutil.RollbackOrLog(ctx, u.logger, err.Error())

		return err
	}

	err = trx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit trx: %w", err)
	}

	return nil
}

func (u *Usecase) processTrx(ctx context.Context, token string) error {
	ver, err := u.emailVerificationsRepository.GetByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("get verification by token: %w", err)
	}

	if ver == nil {
		return dto.ErrInvalidToken
	}

	err = u.userRepository.MarkAsVerified(ctx, ver.UserID)
	if err != nil {
		return fmt.Errorf("mark user as verified: %w", err)
	}

	err = u.emailVerificationsRepository.DeleteTokensByUserID(ctx, ver.UserID)
	if err != nil {
		return fmt.Errorf("delete verification tokens by user id: %w", err)
	}

	return nil
}
