package delete_v1_orgns_orgn_id

import (
	"context"

	app_errors "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/errors"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// DeleteV1OrgnsOrgnId Delete the organization
// (DELETE /v1/orgns/{orgn_id})
func (h *Handler) DeleteV1OrgnsOrgnId(ctx context.Context, req openapi.DeleteV1OrgnsOrgnIdRequestObject) (openapi.DeleteV1OrgnsOrgnIdResponseObject, error) {
	return nil, app_errors.ErrUnimplemented
}
