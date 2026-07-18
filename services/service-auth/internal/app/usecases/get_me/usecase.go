package get_me

import (
	"context"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/orgn"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/get_me/dto"
)

type userRepository interface {
	Find(ctx context.Context, id string) (*user.User, error)
}

type orgnRepository interface {
	Find(ctx context.Context, id string) (*orgn.Orgn, error)
}

type Usecase struct {
	userRepository userRepository
	orgnRepository orgnRepository
}

func NewUsecase(
	userRepository userRepository,
	orgnRepository orgnRepository,
) *Usecase {
	return &Usecase{
		userRepository: userRepository,
		orgnRepository: orgnRepository,
	}
}

func (u *Usecase) Do(ctx context.Context, claims *jwt.Claims) (*dto.Response, error) {
	usr, err := u.userRepository.Find(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	var org *orgn.Orgn
	if claims.OrgnID != nil {
		org, err = u.orgnRepository.Find(ctx, *claims.OrgnID)
		if err != nil {
			return nil, fmt.Errorf("find org: %w", err)
		}
	}

	return &dto.Response{
		SessionID: claims.SessionID,
		User: dto.ResponseUser{
			Email: usr.Email,
			Name:  usr.Name,
		},
		Orgn: &dto.ResponseOrgn{
			ID:   org.ID,
			Name: org.Name,
		},
	}, nil
}
