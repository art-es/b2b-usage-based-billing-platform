package delete_v1_auth_sessions_session_id

import (
	"context"

	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type AuthService interface {
	FinishSession(ctx context.Context, sessionID string) error
}

type Handler struct {
	authService AuthService
}

func NewHandler(authService AuthService) *Handler {
	return &Handler{authService: authService}
}

// DeleteV1AuthSessionsSessionId Finish the session
// (DELETE /v1/auth/sessions/{sessionId})
func (h *Handler) DeleteV1AuthSessionsSessionId(ctx context.Context, req openapi.DeleteV1AuthSessionsSessionIdRequestObject) (openapi.DeleteV1AuthSessionsSessionIdResponseObject, error) {
	err := h.authService.FinishSession(ctx, req.SessionId)
	if err != nil {
		return nil, err
	}

	return &openapi.DeleteV1AuthSessionsSessionId204Response{}, nil
}
