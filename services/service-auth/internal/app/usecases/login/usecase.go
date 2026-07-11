package login

import (
	"context"
	"errors"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/repository"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/services/hash"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/login/dto"
)

type userRepository interface {
	FindByEmail(ctx context.Context, email string) (*user.User, error)
}

type hashService interface {
	Compare(s string, hash string) error
}

type jwtService interface {
	Generate()
}

type Usecase struct {
	userRepository userRepository
	hashService    hashService
}

func (u *Usecase) Do(ctx context.Context, req *dto.Request) (*dto.Response, error) {
	usr, err := u.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, dto.ErrWrongCredentials
		}

		return nil, fmt.Errorf("get user by email: %w", err)
	}

	if !usr.IsVerified {
		return nil, dto.ErrEmailNotVerified
	}

	err = u.hashService.Compare(req.Password, usr.PasswordHash)
	if err != nil {
		if errors.Is(err, hash.ErrMismatch) {
			return nil, dto.ErrWrongCredentials
		}

		return nil, fmt.Errorf("compare password with hash: %w", err)
	}

	return &dto.Response{
		AccessToken:  "",
		RefreshToken: "",
	}, nil
}
