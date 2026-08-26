package put_v1_orgns_orgn_id

import (
	"context"

	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/auth"
	dto "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/clients/orgn_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type authorizer interface {
	Authorize(ctx context.Context) (*auth.Auth, error)
}

type OrgnService interface {
	Update(ctx context.Context, req *dto.UpdateRequest) error
}

type Handler struct {
	authorizer  authorizer
	orgnService OrgnService
}

func NewHandler(
	authorizer authorizer,
	orgnService OrgnService,
) *Handler {
	return &Handler{
		authorizer:  authorizer,
		orgnService: orgnService,
	}
}

// PutV1OrgnsOrgnId Update the organization
// (PUT /v1/orgns/{orgn_id})
func (h *Handler) PutV1OrgnsOrgnId(ctx context.Context, req openapi.PutV1OrgnsOrgnIdRequestObject) (openapi.PutV1OrgnsOrgnIdResponseObject, error) {
	authObj, err := h.authorizer.Authorize(ctx)
	if err != nil {
		return nil, err
	}

	err = h.orgnService.Update(ctx, &dto.UpdateRequest{
		UserID: authObj.UserID,
		OrgnID: req.OrgnId,
		Name:   req.Body.Name,
	})
	if err != nil {
		return nil, err
	}

	return &openapi.PutV1OrgnsOrgnId202Response{}, nil
}
