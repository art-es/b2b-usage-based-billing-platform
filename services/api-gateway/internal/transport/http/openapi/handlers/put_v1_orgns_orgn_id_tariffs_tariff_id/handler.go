package put_v1_orgns_orgn_id_tariffs_tariff_id

import (
	"context"

	app_errors "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/errors"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// PutV1OrgnsOrgnIdTariffsTariffId Update the tariff
// (PUT /v1/orgns/{orgn_id}/tariffs/{tariff_id})
func (h *Handler) PutV1OrgnsOrgnIdTariffsTariffId(ctx context.Context, req openapi.PutV1OrgnsOrgnIdTariffsTariffIdRequestObject) (openapi.PutV1OrgnsOrgnIdTariffsTariffIdResponseObject, error) {
	return nil, app_errors.ErrUnimplemented
}
