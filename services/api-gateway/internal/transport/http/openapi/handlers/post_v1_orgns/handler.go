package post_v1_orgns

import (
	"context"

	app_errors "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/errors"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// PostV1Orgns Create a new organization
// (POST /v1/orgns)
func (h *Handler) PostV1Orgns(ctx context.Context, req openapi.PostV1OrgnsRequestObject) (openapi.PostV1OrgnsResponseObject, error) {
	return nil, app_errors.ErrUnimplemented
}
