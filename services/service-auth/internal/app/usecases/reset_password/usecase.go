package reset_password

import (
	"context"
	"errors"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/repository"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/trx"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/trx/trxutil"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/reset_password/dto"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
)

type passwordResetRepository interface {
	Find(ctx context.Context, token string) (*user.PasswordReset, error)
	DeleteForUser(ctx context.Context, userID string) error
}

type userRepository interface {
	UpdatePasswordHash(ctx context.Context, userID, passwordHash string) error
}

type hashService interface {
	Generate(s string) (string, error)
}

type Usecase struct {
	passwordResetRepository passwordResetRepository
	userRepository          userRepository
	hashService             hashService
	logger                  log.Logger
}

func NewUsecase(
	passwordResetRepository passwordResetRepository,
	userRepository userRepository,
	hashService hashService,
	logger log.Logger,
) *Usecase {
	logger = logger.Set("pkg", "internal/app/usecases/reset_password")

	return &Usecase{
		passwordResetRepository: passwordResetRepository,
		userRepository:          userRepository,
		hashService:             hashService,
		logger:                  logger,
	}
}

func (u *Usecase) Do(ctx context.Context, req *dto.Request) error {
	err := validateRequest(req)
	if err != nil {
		return err
	}

	newPasswordHash, err := u.hashService.Generate(req.NewPassword)
	if err != nil {
		return fmt.Errorf("generate new password hash: %w", err)
	}

	ctx = trx.Begin(ctx)
	err = func() error {
		reset, err := u.passwordResetRepository.Find(ctx, req.Token)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return errInvalidToken
			}

			return fmt.Errorf("find password reset: %w", err)
		}

		err = u.userRepository.UpdatePasswordHash(ctx, reset.UserID, newPasswordHash)
		if err != nil {
			return fmt.Errorf("update password hash: %w", err)
		}

		err = u.passwordResetRepository.DeleteForUser(ctx, reset.UserID)
		if err != nil {
			return fmt.Errorf("delete password resets for user: %w", err)
		}

		return nil
	}()
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
