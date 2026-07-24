package openapi

import (
	"context"

	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/clients/auth_service"
)

type authService interface {
	Login(ctx context.Context, req *auth_service.LoginRequest) (*auth_service.LoginResponse, error)
}
