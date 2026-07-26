package post_v1_auth_register

import (
	"context"

	dto "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/clientdto/auth_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) error
}

type Handler struct {
	authService AuthService
}

func NewHandler(authService AuthService) *Handler {
	return &Handler{authService: authService}
}

func (h *Handler) PostV1AuthRegister(ctx context.Context, req openapi.PostV1AuthRegisterRequestObject) (openapi.PostV1AuthRegisterResponseObject, error) {
	err := h.authService.Register(ctx, &dto.RegisterRequest{
		Name:     req.Body.Name,
		Email:    string(req.Body.Email),
		Password: req.Body.Password,
	})
	if err != nil {
		return nil, err
	}

	return openapi.PostV1AuthRegister202JSONResponse{}, nil
}
