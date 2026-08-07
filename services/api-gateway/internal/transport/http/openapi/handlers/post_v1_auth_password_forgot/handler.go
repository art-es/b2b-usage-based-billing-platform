package post_v1_auth_password_forgot

import (
	"context"

	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type AuthService interface {
	CreatePasswordReset(ctx context.Context, email string) error
}

type Handler struct {
	authService AuthService
}

func NewHandler(authService AuthService) *Handler {
	return &Handler{authService: authService}
}

// PostV1AuthPasswordForgot Send email with reset password link
// (POST /v1/auth/password/forgot)
func (h *Handler) PostV1AuthPasswordForgot(ctx context.Context, req openapi.PostV1AuthPasswordForgotRequestObject) (openapi.PostV1AuthPasswordForgotResponseObject, error) {
	err := h.authService.CreatePasswordReset(ctx, string(req.Body.Email))
	if err != nil {
		return nil, err
	}

	return &openapi.PostV1AuthPasswordForgot202Response{}, nil
}
