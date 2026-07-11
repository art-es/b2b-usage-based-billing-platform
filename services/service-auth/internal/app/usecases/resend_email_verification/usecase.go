package resend_email_verification

import (
	"context"
	"errors"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/resend_email_verification/dto"
)

type emailVerificationRepository interface {
	HasUnsentByEmail(ctx context.Context, email string) (bool, error)
	CreateForEmail(ctx context.Context, email string) error
}

type Usecase struct {
	emailVerificationRepository emailVerificationRepository
}

func NewUsecase(emailVerificationRepository emailVerificationRepository) *Usecase {
	return &Usecase{emailVerificationRepository: emailVerificationRepository}
}

func (u *Usecase) Do(ctx context.Context, email string) error {
	hasUnsent, err := u.emailVerificationRepository.HasUnsentByEmail(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrEmailVerified):
			return dto.ErrEmailVerified
		case errors.Is(err, user.ErrUserNotFound):
			return dto.ErrInvalidEmail
		default:
			return fmt.Errorf("check existence of unsent verification: %w", err)
		}
	}

	if hasUnsent {
		return nil
	}

	err = u.emailVerificationRepository.CreateForEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("create verification for email: %w", err)
	}

	return nil
}
