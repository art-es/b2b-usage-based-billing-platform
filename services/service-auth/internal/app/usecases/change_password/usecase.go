package change_password

import (
	"context"
	"errors"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/hash"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/change_password/dto"
)

type userRepository interface {
	Find(ctx context.Context, id string) (*user.User, error)
	UpdatePasswordHash(ctx context.Context, userID, passwordHash string) error
}

type hashService interface {
	Compare(s string, hash string) error
	Generate(s string) (string, error)
}

type Usecase struct {
	userRepository userRepository
	hashService    hashService
}

func NewUsecase(
	userRepository userRepository,
	hashService hashService,
) *Usecase {
	return &Usecase{
		userRepository: userRepository,
		hashService:    hashService,
	}
}

func (u *Usecase) Do(ctx context.Context, req *dto.Request) error {
	err := validateRequest(req)
	if err != nil {
		return err
	}

	err = u.checkOldPassword(ctx, req.Auth.UserID, req.OldPassword)
	if err != nil {
		return err
	}

	err = u.updatePassword(ctx, req.Auth.UserID, req.NewPassword)
	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) checkOldPassword(ctx context.Context, userID, oldPassword string) error {
	usr, err := u.userRepository.Find(ctx, userID)
	if err != nil {
		return fmt.Errorf("find user: %w", err)
	}

	err = u.hashService.Compare(oldPassword, usr.PasswordHash)
	if err != nil {
		if errors.Is(err, hash.ErrMismatch) {
			return errIncorrectOldPassword
		}

		return fmt.Errorf("compare old password with hash: %w", err)
	}

	return nil
}

func (u *Usecase) updatePassword(ctx context.Context, userID, newPassword string) error {
	newPasswordHash, err := u.hashService.Generate(newPassword)
	if err != nil {
		return fmt.Errorf("generate new password hash: %w", err)
	}

	err = u.userRepository.UpdatePasswordHash(ctx, userID, newPasswordHash)
	if err != nil {
		return fmt.Errorf("update password hash: %w", err)
	}

	return nil
}
