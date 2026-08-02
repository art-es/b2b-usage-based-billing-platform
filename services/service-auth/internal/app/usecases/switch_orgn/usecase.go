package switch_orgn

import (
	"context"
	"errors"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/orgn"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/repository"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/switch_orgn/dto"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/validate"
)

type sessionRepository interface {
	SetOrgnID(ctx context.Context, sessID, orgnID string) error
}

type orgnService interface {
	GetByID(ctx context.Context, orgnID string) (*orgn.Orgn, error)
}

type Usecase struct {
	sessionRepository sessionRepository
	orgnService       orgnService
}

func NewUsecase(
	sessionRepository sessionRepository,
	orgnService orgnService,
) *Usecase {
	return &Usecase{
		sessionRepository: sessionRepository,
		orgnService:       orgnService,
	}
}

func (u *Usecase) Do(ctx context.Context, req *dto.Request) error {
	org, err := u.orgnService.GetByID(ctx, req.OrgnID)
	if err != nil {
		if errors.Is(err, orgn.ErrNotFound) {
			return errIncorrectOrgnID
		}

		return fmt.Errorf("get orgn by id: %w", err)
	}

	if org.UserID != req.Auth.UserID {
		return errIncorrectOrgnID
	}

	err = u.sessionRepository.SetOrgnID(ctx, req.Auth.SessionID, req.Auth.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return validate.ErrUnauthorized
		}

		return fmt.Errorf("set orgn id to session: %w", err)
	}

	return nil
}
