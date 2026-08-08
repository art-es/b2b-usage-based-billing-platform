package get_v1_webhook

import (
	"context"

	app_errors "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/errors"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// GetV1Webhook Connect to webhook to receive events
// (GET /v1/webhook)
func (h *Handler) GetV1Webhook(ctx context.Context, req openapi.GetV1WebhookRequestObject) (openapi.GetV1WebhookResponseObject, error) {
	return nil, app_errors.ErrUnimplemented
}
