//go:generate mockgen -source=usecase.go -destination=usecase_mock_test.go -package=$GOPACKAGE
package get_me

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/orgn"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/get_me/dto"
)

type userRepository interface {
	Find(ctx context.Context, id string) (*user.User, error)
}

type orgnService interface {
	GetByID(ctx context.Context, id string) (*orgn.Orgn, error)
}

type Usecase struct {
	userRepository userRepository
	orgnService    orgnService
}

func NewUsecase(
	userRepository userRepository,
	orgnService orgnService,
) *Usecase {
	return &Usecase{
		userRepository: userRepository,
		orgnService:    orgnService,
	}
}

func (u *Usecase) Do(ctx context.Context, claims *jwt.Claims) (*dto.Response, error) {
	var usr *user.User
	var org *orgn.Orgn

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() (err error) {
		usr, err = u.userRepository.Find(egCtx, claims.UserID)
		if err != nil {
			return fmt.Errorf("find user: %w", err)
		}

		return nil
	})

	if claims.OrgnID != nil {
		eg.Go(func() (err error) {
			org, err = u.orgnService.GetByID(egCtx, *claims.OrgnID)
			if err != nil {
				return fmt.Errorf("get orgn by id: %w", err)
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return &dto.Response{
		SessionID: claims.SessionID,
		User: dto.ResponseUser{
			Email: usr.Email,
			Name:  usr.Name,
		},
		Orgn: convertOrgn(org),
	}, nil
}

func convertOrgn(in *orgn.Orgn) *dto.ResponseOrgn {
	if in != nil {
		return &dto.ResponseOrgn{
			ID:   in.ID,
			Name: in.Name,
		}
	}
	return nil
}
