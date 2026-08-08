package get_v1_orgns_orgn_id_api_keys

import (
	"context"

	app_errors "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/errors"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// GetV1OrgnsOrgnIdApiKeys Get orgn API keys
// (GET /v1/orgns/{orgn_id}/api-keys)
func (h *Handler) GetV1OrgnsOrgnIdApiKeys(ctx context.Context, req openapi.GetV1OrgnsOrgnIdApiKeysRequestObject) (openapi.GetV1OrgnsOrgnIdApiKeysResponseObject, error) {
	return nil, app_errors.ErrUnimplemented
}
