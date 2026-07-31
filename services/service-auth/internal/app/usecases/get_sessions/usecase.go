package get_sessions

import (
	"context"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/session"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/get_sessions/dto"
)

type sessionsRepository interface {
	Get(ctx context.Context, userID string, cursor *session.ListCursor) (list []*session.Session, nextCursor *session.ListCursor, err error)
}

type keyedHashService interface {
	Generate(secret []byte, s string) (string, error)
}

type Usecase struct {
	sessionsRepository sessionsRepository
	keyedHashService   keyedHashService
	cursorSecretKey    []byte
}

func NewUsecase(
	sessionsRepository sessionsRepository,
	keyedHashService keyedHashService,
	cursorSecretKey string,
) *Usecase {
	return &Usecase{
		sessionsRepository: sessionsRepository,
		keyedHashService:   keyedHashService,
		cursorSecretKey:    []byte(cursorSecretKey),
	}
}

func (u *Usecase) Do(ctx context.Context, req *dto.Request) (*dto.Response, error) {
	currCursor, err := u.stringToCursor(req.Cursor)
	if err != nil {
		return nil, err
	}

	list, nextCursor, err := u.sessionsRepository.Get(ctx, req.Auth.UserID, currCursor)
	if err != nil {
		return nil, fmt.Errorf("get sessions from repository: %w", err)
	}

	nextCursorHash, err := u.cursorToString(nextCursor)
	if err != nil {
		return nil, err
	}

	return &dto.Response{
		Sessions:   convertSessions(list),
		NextCursor: nextCursorHash,
	}, nil
}

func convertSessions(in []*session.Session) []*dto.Session {
	out := make([]*dto.Session, 0, len(in))
	for _, s := range in {
		out = append(out, convertSession(s))
	}
	return out
}

func convertSession(in *session.Session) *dto.Session {
	return &dto.Session{
		ID:        in.ID,
		CreatedAt: in.CreatedAt,
	}
}
