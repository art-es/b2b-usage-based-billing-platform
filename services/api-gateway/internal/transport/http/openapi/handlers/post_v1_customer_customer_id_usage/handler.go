package post_v1_customer_customer_id_usage

import (
	"context"

	app_errors "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/errors"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// PostV1CustomerCustomerIdUsage Add new usage
// (POST /v1/customer/{customer_id}/usage)
func (h *Handler) PostV1CustomerCustomerIdUsage(ctx context.Context, req openapi.PostV1CustomerCustomerIdUsageRequestObject) (openapi.PostV1CustomerCustomerIdUsageResponseObject, error) {
	return nil, app_errors.ErrUnimplemented
}
