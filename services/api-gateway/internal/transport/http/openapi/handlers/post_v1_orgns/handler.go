package post_v1_orgns

import (
	"context"

	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/auth"
	dto "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/clients/orgn_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/openapiutil"
)

type authorizer interface {
	Authorize(ctx context.Context) (*auth.Auth, error)
}

type OrgnService interface {
	Create(ctx context.Context, req *dto.CreateRequest) (*dto.CreateResponse, error)
}

type AuthService interface {
	SwitchOrgn(ctx context.Context, orgnID string) (string, error)
}

type Handler struct {
	authorizer  authorizer
	orgnService OrgnService
	authService AuthService
}

func NewHandler(
	authorizer authorizer,
	orgnService OrgnService,
	authService AuthService,
) *Handler {
	return &Handler{
		authorizer:  authorizer,
		orgnService: orgnService,
		authService: authService,
	}
}

// PostV1Orgns Create a new organization
// (POST /v1/orgns)
func (h *Handler) PostV1Orgns(ctx context.Context, req openapi.PostV1OrgnsRequestObject) (openapi.PostV1OrgnsResponseObject, error) {
	authObj, err := h.authorizer.Authorize(ctx)
	if err != nil {
		return nil, err
	}

	createOrgnRes, err := h.orgnService.Create(ctx, &dto.CreateRequest{
		UserID: authObj.UserID,
		Name:   req.Body.Name,
	})
	if err != nil {
		return nil, err
	}

	// TODO: add circuit breaker
	newAccessToken, err := h.authService.SwitchOrgn(ctx, createOrgnRes.OrgnID)
	if err != nil {
		return nil, err
	}

	return openapi.PostV1Orgns201JSONResponse{
		OrgnId:      openapiutil.ToUUID(createOrgnRes.OrgnID),
		AccessToken: newAccessToken,
	}, nil
}
