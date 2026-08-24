package get_sessions

import (
	"context"
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/session"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/get_sessions/dto"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/listcursor"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
)

type sessionsRepository interface {
	Get(ctx context.Context, userID string, cursor *session.ListCursor) (list []*session.Session, nextCursor *session.ListCursor, err error)
}

type Usecase struct {
	sessionsRepository sessionsRepository
	listCursorEncoder  *listcursor.Encoder[session.ListCursor]
	logger             log.Logger
}

func NewUsecase(
	sessionsRepository sessionsRepository,
	cursorSecret string,
	logger log.Logger,
) *Usecase {
	logger = logger.Set("pkg", "internal/app/usecases/get_sessions")

	return &Usecase{
		sessionsRepository: sessionsRepository,
		listCursorEncoder:  listcursor.NewEncoder[session.ListCursor]([]byte(cursorSecret)),
		logger:             logger,
	}
}

func (u *Usecase) Do(ctx context.Context, req *dto.Request) (*dto.Response, error) {
	cursor, err := u.listCursorEncoder.DecodeAndCompare(req.Cursor)
	if err != nil {
		u.logger.Log(log.Warning).
			Set("message", "decode cursor error").
			Set("error", err.Error()).
			Write()
	}

	list, nextCursor, err := u.sessionsRepository.Get(ctx, req.Auth.UserID, cursor)
	if err != nil {
		return nil, fmt.Errorf("get sessions from repository: %w", err)
	}

	res := &dto.Response{
		Sessions: convertSessions(list),
	}

	res.NextCursor, err = u.listCursorEncoder.Encode(nextCursor)
	if err != nil {
		return nil, fmt.Errorf("encode next cursor: %w", err)
	}

	return res, nil
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
