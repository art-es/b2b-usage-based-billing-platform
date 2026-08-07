package post_v1_auth_password_reset

import (
	"context"

	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type AuthService interface {
	ResetPassword(ctx context.Context, token, newPassword string) error
}

type Handler struct {
	authService AuthService
}

func NewHandler(authService AuthService) *Handler {
	return &Handler{authService: authService}
}

// PostV1AuthPasswordReset Reset the password
// (POST /v1/auth/password/reset)
func (h *Handler) PostV1AuthPasswordReset(ctx context.Context, req openapi.PostV1AuthPasswordResetRequestObject) (openapi.PostV1AuthPasswordResetResponseObject, error) {
	err := h.authService.ResetPassword(ctx, req.Body.Token.String(), req.Body.NewPassword)
	if err != nil {
		return nil, err
	}

	return &openapi.PostV1AuthPasswordReset204Response{}, nil
}
