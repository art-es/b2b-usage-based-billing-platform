package finish_all_sessions

import (
	"context"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
)

type sessionsRepository interface {
	DeleteAll(ctx context.Context, userID string) error
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

func (u *Usecase) Do(ctx context.Context, auth *jwt.Claims) error {
	err := u.sessionsRepository.DeleteAll(ctx, auth.UserID)
	if err != nil {
		return fmt.Errorf("delete all sessions: %w", err)
	}

	return nil
}
