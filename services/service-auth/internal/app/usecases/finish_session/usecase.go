package finish_session

import (
	"context"
	"errors"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/repository"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/finish_session/dto"
)

type sessionsRepository interface {
	Delete(ctx context.Context, userID, sessionID string) error
}

type Usecase struct {
	sessionsRepository sessionsRepository
}

func NewUsecase(
	sessionsRepository sessionsRepository,
) *Usecase {
	return &Usecase{
		sessionsRepository: sessionsRepository,
	}
}

func (u *Usecase) Do(ctx context.Context, req *dto.Request) error {
	err := u.sessionsRepository.Delete(ctx, req.Auth.UserID, req.SessionID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return errIncorrectSessionID
		}

		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}
