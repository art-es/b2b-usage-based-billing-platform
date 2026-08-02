package delete_v1_auth_sessions

import (
	"context"

	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type AuthService interface {
	FinishAllSessions(ctx context.Context) error
}

type Handler struct {
	authService AuthService
}

func NewHandler(authService AuthService) *Handler {
	return &Handler{authService: authService}
}

// DeleteV1AuthSessions Finish all sessions
// (DELETE /v1/auth/sessions)
func (h *Handler) DeleteV1AuthSessions(ctx context.Context, req openapi.DeleteV1AuthSessionsRequestObject) (openapi.DeleteV1AuthSessionsResponseObject, error) {
	err := h.authService.FinishAllSessions(ctx)
	if err != nil {
		return nil, err
	}

	return &openapi.DeleteV1AuthSessions204Response{}, nil
}
