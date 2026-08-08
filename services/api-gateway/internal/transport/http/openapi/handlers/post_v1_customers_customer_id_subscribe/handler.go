package post_v1_customers_customer_id_subscribe

import (
	"context"

	app_errors "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/errors"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// PostV1CustomersCustomerIdSubscribe Subscribe customer to tariff
// (POST /v1/customers/{customer_id}/subscribe)
func (h *Handler) PostV1CustomersCustomerIdSubscribe(ctx context.Context, req openapi.PostV1CustomersCustomerIdSubscribeRequestObject) (openapi.PostV1CustomersCustomerIdSubscribeResponseObject, error) {
	return nil, app_errors.ErrUnimplemented
}
