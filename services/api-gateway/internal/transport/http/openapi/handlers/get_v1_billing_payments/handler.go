package get_v1_billing_payments

import (
	"context"

	app_errors "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/errors"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// GetV1BillingPayments Get payment status
// (GET /v1/billing/payments)
func (h *Handler) GetV1BillingPayments(ctx context.Context, req openapi.GetV1BillingPaymentsRequestObject) (openapi.GetV1BillingPaymentsResponseObject, error) {
	return nil, app_errors.ErrUnimplemented
}
