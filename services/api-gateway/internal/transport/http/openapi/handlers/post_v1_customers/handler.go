package post_v1_customers

import (
	"context"

	app_errors "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/errors"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// PostV1Customers Create a new customer
// (POST /v1/customers)
func (h *Handler) PostV1Customers(ctx context.Context, req openapi.PostV1CustomersRequestObject) (openapi.PostV1CustomersResponseObject, error) {
	return nil, app_errors.ErrUnimplemented
}
