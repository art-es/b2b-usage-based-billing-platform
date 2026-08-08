package get_v1_orgns_orgn_id_members

import (
	"context"

	app_errors "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/errors"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// GetV1OrgnsOrgnIdMembers Get organization members
// (GET /v1/orgns/{orgn_id}/members)
func (h *Handler) GetV1OrgnsOrgnIdMembers(ctx context.Context, req openapi.GetV1OrgnsOrgnIdMembersRequestObject) (openapi.GetV1OrgnsOrgnIdMembersResponseObject, error) {
	return nil, app_errors.ErrUnimplemented
}
