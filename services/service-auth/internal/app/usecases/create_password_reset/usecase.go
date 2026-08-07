package create_password_reset

import (
	"context"
	"errors"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/repository"
)

type userRepository interface {
	FindByEmail(ctx context.Context, email string) (*user.User, error)
}

type passwordResetRepository interface {
	Create(ctx context.Context, userID string) error
}

type Usecase struct {
	userRepository          userRepository
	passwordResetRepository passwordResetRepository
}

func NewUsecase(
	userRepository          userRepository,
	passwordResetRepository passwordResetRepository,
) *Usecase {
	return &Usecase{
		userRepository:          userRepository,
		passwordResetRepository: passwordResetRepository,
	}
}

func (u *Usecase) Do(ctx context.Context, email string) error {
	err := validateEmail(email)
	if err != nil {
		return err
	}

	usr, err := u.userRepository.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return errIncorrectEmail
		}

		return fmt.Errorf("find by user: %w", err)
	}

	err = u.passwordResetRepository.Create(ctx, usr.ID)
	if err != nil {
		return fmt.Errorf("create password reset: %w", err)
	}

	return nil
}
