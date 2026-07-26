package post_v1_auth_email_resend_verification

import (
	"context"

	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type AuthService interface {
	ResendEmailVerification(ctx context.Context, email string) error
}

type Handler struct {
	authService AuthService
}

func NewHandler(authService AuthService) *Handler {
	return &Handler{authService: authService}
}

// PostV1AuthEmailResendVerification Resend email verification
// (POST /v1/auth/email/resend-verification)
func (h *Handler) PostV1AuthEmailResendVerification(ctx context.Context, req openapi.PostV1AuthEmailResendVerificationRequestObject) (openapi.PostV1AuthEmailResendVerificationResponseObject, error) {
	err := h.authService.ResendEmailVerification(ctx, req.Body.Email)
	if err != nil {
		return nil, err
	}

	return &openapi.PostV1AuthEmailResendVerification204Response{}, nil
}
