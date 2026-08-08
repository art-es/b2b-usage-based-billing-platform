package get_v1_orgns_orgn_id_tariffs

import (
	"context"

	app_errors "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/errors"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// GetV1OrgnsOrgnIdTariffs Get tariffs
// (GET /v1/orgns/{orgn_id}/tariffs)
func (h *Handler) GetV1OrgnsOrgnIdTariffs(ctx context.Context, req openapi.GetV1OrgnsOrgnIdTariffsRequestObject) (openapi.GetV1OrgnsOrgnIdTariffsResponseObject, error) {
	return nil, app_errors.ErrUnimplemented
}
