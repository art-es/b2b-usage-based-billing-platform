package get

import (
	"context"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/domains/orgn"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/usecases/get/dto"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/pkg/log"
)

type orgnRepository interface {
	Get(ctx context.Context, userID string, cursor *orgn.ListCursor) ([]*orgn.Orgn, *orgn.ListCursor, error)
}

type Usecase struct {
	orgnRepository orgnRepository
	cursorSecret   []byte
	logger         log.Logger
}

func NewUsecase(
	orgnRepository orgnRepository,
	cursorSecret string,
	logger log.Logger,
) *Usecase {
	return &Usecase{
		orgnRepository: orgnRepository,
		cursorSecret:   []byte(cursorSecret),
		logger:         logger,
	}
}

func (u *Usecase) Do(ctx context.Context, req *dto.Request) (*dto.Response, error) {
	cursor, err := orgn.NewListCursorFromString(u.cursorSecret, req.Cursor)
	if err != nil {
		u.logger.Log(log.Warning).
			Set("message", "make cursor from string error").
			Set("error", err.Error()).
			Write()
	}

	list, nextCursor, err := u.orgnRepository.Get(ctx, req.UserID, cursor)
	if err != nil {
		return nil, fmt.Errorf("get orgn: %w", err)
	}

	nextCursorStr, err := nextCursor.String(u.cursorSecret)
	if err != nil {
		return nil, fmt.Errorf("convert cursor to string: %w", err)
	}

	return &dto.Response{
		List:       list,
		NextCursor: nextCursorStr,
	}, nil
}
