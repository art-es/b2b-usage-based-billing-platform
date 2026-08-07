package post_v1_auth_password_change

import (
	"context"

	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type AuthService interface {
	ChangePassword(ctx context.Context, oldPassword, newPassword string) error
}

type Handler struct {
	authService AuthService
}

func NewHandler(authService AuthService) *Handler {
	return &Handler{authService: authService}
}

// PostV1AuthPasswordChange Change the password of authorized user
// (POST /v1/auth/password/change)
func (h *Handler) PostV1AuthPasswordChange(ctx context.Context, req openapi.PostV1AuthPasswordChangeRequestObject) (openapi.PostV1AuthPasswordChangeResponseObject, error) {
	err := h.authService.ChangePassword(ctx, req.Body.OldPassword, req.Body.NewPassword)
	if err != nil {
		return nil, err
	}

	return &openapi.PostV1AuthPasswordChange204Response{}, nil
}
