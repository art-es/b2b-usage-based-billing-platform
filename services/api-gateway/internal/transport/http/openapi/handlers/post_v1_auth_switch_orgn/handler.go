package post_v1_auth_switch_orgn

import (
	"context"

	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/openapiutil"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type AuthService interface {
	SwitchOrgn(ctx context.Context, orgnID string) error
}

type Handler struct {
	authService AuthService
}

func NewHandler(authService AuthService) *Handler {
	return &Handler{authService: authService}
}

// GetV1Me Get user info of current session
// (GET /v1/me)
func (h *Handler) GetV1Me(ctx context.Context, req openapi.GetV1MeRequestObject) (openapi.GetV1MeResponseObject, error) {
	dtoRes, err := h.authService.GetMe(ctx)
	if err != nil {
		return nil, err
	}

	res := openapi.GetV1Me200JSONResponse{
		SessionId: dtoRes.SessionID,
		Name:      dtoRes.Name,
		Email:     openapi_types.Email(dtoRes.Email),
	}

	if dtoRes.Orgn != nil {
		res.Orgn = &openapi.MeResponseOrgn{
			Id:   openapiutil.ToUUID(dtoRes.Orgn.ID),
			Name: dtoRes.Name,
		}
	}

	return res, nil
}
