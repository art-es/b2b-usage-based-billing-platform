package get_v1_auth_sessions

import (
	"context"
	"time"

	dto "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/clients/auth_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type AuthService interface {
	GetSessions(ctx context.Context, req *dto.GetSessionsRequest) (*dto.GetSessionsResponse, error)
}

type Handler struct {
	authService AuthService
}

func NewHandler(authService AuthService) *Handler {
	return &Handler{authService: authService}
}

// GetV1AuthSessions Get sessions
// (GET /v1/auth/sessions)
func (h *Handler) GetV1AuthSessions(ctx context.Context, req openapi.GetV1AuthSessionsRequestObject) (openapi.GetV1AuthSessionsResponseObject, error) {
	dtoRes, err := h.authService.GetSessions(ctx, &dto.GetSessionsRequest{
		Cursor: req.Body.Cursor,
	})
	if err != nil {
		return nil, err
	}

	sessions := make([]openapi.Session, 0, len(dtoRes.Sessions))
	for _, s := range dtoRes.Sessions {
		sessions = append(sessions, openapi.Session{
			Id:        s.ID,
			CreatedAt: s.CreatedAt.Format(time.RFC3339),
		})
	}

	return &openapi.GetV1AuthSessions200JSONResponse{
		Sessions:   sessions,
		NextCursor: dtoRes.NextCursor,
	}, nil
}
