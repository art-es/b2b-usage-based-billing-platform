package patch_v1_orgns_orgn_id_members_member_id

import (
	"context"

	app_errors "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/errors"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// PatchV1OrgnsOrgnIdMembersMemberId Update organization member
// (PATCH /v1/orgns/{orgn_id}/members/{member_id})
func (h *Handler) PatchV1OrgnsOrgnIdMembersMemberId(ctx context.Context, req openapi.PatchV1OrgnsOrgnIdMembersMemberIdRequestObject) (openapi.PatchV1OrgnsOrgnIdMembersMemberIdResponseObject, error) {
	return nil, app_errors.ErrUnimplemented
}
