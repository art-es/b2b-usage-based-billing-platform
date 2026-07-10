//go:generate mockgen -source=usecase.go -destination=usecase_mock_test.go -package=$GOPACKAGE
package register

import (
	"context"
	"errors"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/repository"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/register/dto"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/trx"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/trx/trxutil"
)

type userRepository interface {
	Create(ctx context.Context, user *user.User) error
}

type verificationRepository interface {
	Create(ctx context.Context, userID string) error
}

type hashService interface {
	Generate(s string) (string, error)
}

type Usecase struct {
	hashService            hashService
	userRepository         userRepository
	verificationRepository verificationRepository
	logger                 log.Logger
}

func NewUsecase(
	hashService hashService,
	userRepository userRepository,
	verificationRepository verificationRepository,
	logger log.Logger,
) *Usecase {
	logger = logger.Set("pkg", "internal/app/usecases/register")

	return &Usecase{
		hashService:            hashService,
		userRepository:         userRepository,
		verificationRepository: verificationRepository,
		logger:                 logger,
	}
}

func (u *Usecase) Do(ctx context.Context, req *dto.Request) error {
	passwordHash, err := u.hashService.Generate(req.Password)
	if err != nil {
		return fmt.Errorf("generate password hash: %w", err)
	}

	usr := user.NewRegisteredUser(req.Name, req.Email, passwordHash)

	ctx = trx.Begin(ctx)

	err = u.userRepository.Create(ctx, usr)
	if err != nil {
		trxutil.RollbackOrLog(ctx, u.logger, fmt.Sprintf("create user: %v", err))

		if errors.Is(err, repository.ErrUnique) {
			return dto.ErrEmailInUse
		}

		return fmt.Errorf("create user: %w", err)
	}

	err = u.verificationRepository.Create(ctx, usr.ID)
	if err != nil {
		trxutil.RollbackOrLog(ctx, u.logger, fmt.Sprintf("create verification: %v", err))

		return fmt.Errorf("create verification: %w", err)
	}

	err = trx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit trx: %w", err)
	}

	return nil
}
