package post_v1_auth_email_verify

import (
	"context"

	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type AuthService interface {
	VerifyEmail(ctx context.Context, token string) error
}

type Handler struct {
	authService AuthService
}

func NewHandler(authService AuthService) *Handler {
	return &Handler{authService: authService}
}

// PostV1AuthEmailVerify Verify email
// (POST /v1/auth/email/verify)
func (h *Handler) PostV1AuthEmailVerify(ctx context.Context, req openapi.PostV1AuthEmailVerifyRequestObject) (openapi.PostV1AuthEmailVerifyResponseObject, error) {
	err := h.authService.VerifyEmail(ctx, req.Body.Token.String())
	if err != nil {
		return nil, err
	}

	return &openapi.PostV1AuthEmailVerify204Response{}, nil
}
