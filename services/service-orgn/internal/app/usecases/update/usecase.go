package update

import (
	"context"
	"errors"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/domains/orgn"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/repository"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/trx"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/trx/trxutil"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/usecases/update/dto"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/pkg/log"
)

type orgnRepository interface {
	Find(ctx context.Context, orgnID, userID string) (*orgn.Orgn, error)
	Save(ctx context.Context, org *orgn.Orgn) error
}

type Usecase struct {
	orgnRepository orgnRepository
	logger         log.Logger
}

func NewUsecase(orgnRepository orgnRepository) *Usecase {
	return &Usecase{
		orgnRepository: orgnRepository,
	}
}

func (u *Usecase) Do(ctx context.Context, req *dto.Request) error {
	ctx = trx.Begin(ctx)
	err := func() error {
		org, err := u.orgnRepository.Find(ctx, req.OrgnID, req.UserID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return errOrgnIDNotFound
			}

			return fmt.Errorf("find orgn: %w", err)
		}

		if hasUpdates := org.Update(req.Name); !hasUpdates {
			return nil
		}

		err = u.orgnRepository.Save(ctx, org)
		if err != nil {
			return fmt.Errorf("save orgn: %w", err)
		}

		return nil
	}()
	if err != nil {
		trxutil.RollbackOrLog(ctx, u.logger, err.Error())
		return err
	}

	err = trx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("trx commit: %w", err)
	}

	return nil
}
