package create

import (
	"context"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/domains/orgn"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/usecases/create/dto"
)

type orgnRepository interface {
	Save(ctx context.Context, org *orgn.Orgn) error
}

type Usecase struct {
	orgnRepository orgnRepository
}

func NewUsecase(orgnRepository orgnRepository) *Usecase {
	return &Usecase{
		orgnRepository: orgnRepository,
	}
}

func (u *Usecase) Do(ctx context.Context, req *dto.Request) (*dto.Response, error) {
	err := validateRequest(req)
	if err != nil {
		return nil, err
	}

	org := orgn.New(req.Name, req.UserID)
	err = u.orgnRepository.Save(ctx, org)
	if err != nil {
		return nil, fmt.Errorf("save orgn: %w", err)
	}

	return &dto.Response{
		OrgnID: org.ID,
	}, nil
}
