package post_v1_orgns_orgn_id_api_keys

import (
	"context"

	app_errors "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/errors"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// PostV1OrgnsOrgnIdApiKeys Create a new organization API key
// (POST /v1/orgns/{orgn_id}/api-keys)
func (h *Handler) PostV1OrgnsOrgnIdApiKeys(ctx context.Context, req openapi.PostV1OrgnsOrgnIdApiKeysRequestObject) (openapi.PostV1OrgnsOrgnIdApiKeysResponseObject, error) {
	return nil, app_errors.ErrUnimplemented
}
