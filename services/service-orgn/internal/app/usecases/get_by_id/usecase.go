package get_by_id

import (
	"context"
	"errors"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/domains/orgn"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/repository"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/usecases/get_by_id/dto"
)

type orgnRepository interface {
	Find(ctx context.Context, orgnID, userID string) (*orgn.Orgn, error)
}

type Usecase struct {
	orgnRepository orgnRepository
}

func NewUsecase(orgnRepository orgnRepository) *Usecase {
	return &Usecase{
		orgnRepository: orgnRepository,
	}
}

func (u *Usecase) Do(ctx context.Context, req *dto.Request) (*orgn.Orgn, error) {
	err := validateRequest(req)
	if err != nil {
		return nil, err
	}

	org, err := u.orgnRepository.Find(ctx, req.OrgnID, req.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, errOrgnIDNotFound
		}

		return nil, fmt.Errorf("find orgn: %w", err)
	}

	return org, nil
}
