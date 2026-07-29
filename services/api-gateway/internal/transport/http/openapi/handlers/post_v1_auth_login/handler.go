package post_v1_auth_login

import (
	"context"

	dto "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/clients/auth_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type AuthService interface {
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
}

type Handler struct {
	authService AuthService
}

func NewHandler(authService AuthService) *Handler {
	return &Handler{authService: authService}
}

// PostV1AuthLogin Create a new session
// (POST /v1/auth/login)
func (h *Handler) PostV1AuthLogin(ctx context.Context, req openapi.PostV1AuthLoginRequestObject) (openapi.PostV1AuthLoginResponseObject, error) {
	res, err := h.authService.Login(ctx, &dto.LoginRequest{
		Email:    string(req.Body.Email),
		Password: req.Body.Password,
	})
	if err != nil {
		return nil, err
	}

	return &openapi.PostV1AuthLogin200JSONResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}
