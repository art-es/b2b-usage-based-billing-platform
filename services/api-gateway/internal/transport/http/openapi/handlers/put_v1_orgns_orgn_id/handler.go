package put_v1_orgns_orgn_id

import (
	"context"

	app_errors "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/errors"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// PutV1OrgnsOrgnId Update the organization
// (PUT /v1/orgns/{orgn_id})
func (h *Handler) PutV1OrgnsOrgnId(ctx context.Context, req openapi.PutV1OrgnsOrgnIdRequestObject) (openapi.PutV1OrgnsOrgnIdResponseObject, error) {
	return nil, app_errors.ErrUnimplemented
}
