package get_v1_orgns

import (
	"context"

	app_errors "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/errors"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// GetV1Orgns Get organizations
// (GET /v1/orgns)
func (h *Handler) GetV1Orgns(ctx context.Context, req openapi.GetV1OrgnsRequestObject) (openapi.GetV1OrgnsResponseObject, error) {
	return nil, app_errors.ErrUnimplemented
}
