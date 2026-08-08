package post_v1_billing_payments

import (
	"context"

	app_errors "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/errors"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// PostV1BillingPayments Generate a new payment link
// (POST /v1/billing/payments)
func (h *Handler) PostV1BillingPayments(ctx context.Context, req openapi.PostV1BillingPaymentsRequestObject) (openapi.PostV1BillingPaymentsResponseObject, error) {
	return nil, app_errors.ErrUnimplemented
}
