package post_v1_auth_switch_orgn

import (
	"context"

	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type AuthService interface {
	SwitchOrgn(ctx context.Context, orgnID string) (string, error)
}

type Handler struct {
	authService AuthService
}

func NewHandler(authService AuthService) *Handler {
	return &Handler{authService: authService}
}

// PostV1AuthSwitchOrgn Switch an organization in session
// (POST /v1/auth/switch-orgn)
func (h *Handler) PostV1AuthSwitchOrgn(ctx context.Context, req openapi.PostV1AuthSwitchOrgnRequestObject) (openapi.PostV1AuthSwitchOrgnResponseObject, error) {
	accessToken, err := h.authService.SwitchOrgn(ctx, req.Body.OrgnId.String())
	if err != nil {
		return nil, err
	}

	return &openapi.PostV1AuthSwitchOrgn200JSONResponse{
		AccessToken: accessToken,
	}, nil
}
