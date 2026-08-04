package get_by_id

import (
	"context"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/domains/orgn"
)

type orgnRepository interface {
	Find(ctx context.Context, orgnID string) (*orgn.Orgn, error)
}

type Usecase struct {
	orgnRepository orgnRepository
}

func NewUsecase(orgnRepository orgnRepository) *Usecase {
	return &Usecase{
		orgnRepository: orgnRepository,
	}
}

func (u *Usecase) Do(ctx context.Context, orgnID string) (*orgn.Orgn, error) {
	return u.orgnRepository.Find(ctx, orgnID)
}
