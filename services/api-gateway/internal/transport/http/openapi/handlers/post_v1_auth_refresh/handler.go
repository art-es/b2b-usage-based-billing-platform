package post_v1_auth_refresh

import (
	"context"

	dto "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/clientdto/auth_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type AuthService interface {
	RefreshSession(ctx context.Context, token string) (*dto.RefreshSessionResponse, error)
}

type Handler struct {
	authService AuthService
}

func NewHandler(authService AuthService) *Handler {
	return &Handler{authService: authService}
}

// PostV1AuthRefresh Refresh an access token of session
// (POST /v1/auth/refresh)
func (h *Handler) PostV1AuthRefresh(ctx context.Context, req openapi.PostV1AuthRefreshRequestObject) (openapi.PostV1AuthRefreshResponseObject, error) {
	res, err := h.authService.RefreshSession(ctx, req.Body.Token)
	if err != nil {
		return nil, err
	}

	return &openapi.PostV1AuthRefresh200JSONResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}
