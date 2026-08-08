package delete_v1_orgns_orgn_id_tariffs_tariff_id

import (
	"context"

	app_errors "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/errors"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// DeleteV1OrgnsOrgnIdTariffsTariffId Delete the tariff
// (DELETE /v1/orgns/{orgn_id}/tariffs/{tariff_id})
func (h *Handler) DeleteV1OrgnsOrgnIdTariffsTariffId(ctx context.Context, req openapi.DeleteV1OrgnsOrgnIdTariffsTariffIdRequestObject) (openapi.DeleteV1OrgnsOrgnIdTariffsTariffIdResponseObject, error) {
	return nil, app_errors.ErrUnimplemented
}
