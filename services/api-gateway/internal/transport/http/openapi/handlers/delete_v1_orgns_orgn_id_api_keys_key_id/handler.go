package delete_v1_orgns_orgn_id_api_keys_key_id

import (
	"context"

	app_errors "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/errors"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// DeleteV1OrgnsOrgnIdApiKeysKeyId Delete organization API key
// (DELETE /v1/orgns/{orgn_id}/api-keys/{key_id})
func (h *Handler) DeleteV1OrgnsOrgnIdApiKeysKeyId(ctx context.Context, req openapi.DeleteV1OrgnsOrgnIdApiKeysKeyIdRequestObject) (openapi.DeleteV1OrgnsOrgnIdApiKeysKeyIdResponseObject, error) {
	return nil, app_errors.ErrUnimplemented
}
