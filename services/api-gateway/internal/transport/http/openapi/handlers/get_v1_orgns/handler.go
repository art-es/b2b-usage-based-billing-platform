package get_v1_orgns

import (
	"context"

	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/auth"
	dto "github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/clients/orgn_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/openapiutil"
)

type authorizer interface {
	Authorize(ctx context.Context) (*auth.Auth, error)
}

type OrgnService interface {
	Get(ctx context.Context, req *dto.GetRequest) (*dto.GetResponse, error)
}

type Handler struct {
	authorizer  authorizer
	orgnService OrgnService
}

func NewHandler(authorizer authorizer, orgnService OrgnService) *Handler {
	return &Handler{
		authorizer:  authorizer,
		orgnService: orgnService,
	}
}

// GetV1Orgns Get organizations
// (GET /v1/orgns)
func (h *Handler) GetV1Orgns(ctx context.Context, req openapi.GetV1OrgnsRequestObject) (openapi.GetV1OrgnsResponseObject, error) {
	auth, err := h.authorizer.Authorize(ctx)
	if err != nil {
		return nil, err
	}

	res, err := h.orgnService.Get(ctx, &dto.GetRequest{
		UserID: auth.UserID,
		Cursor: req.Params.Cursor,
	})
	if err != nil {
		return nil, err
	}

	orgns := make([]openapi.Orgn, 0, len(res.Orgns))
	for _, orgn := range res.Orgns {
		orgns = append(orgns, openapi.Orgn{
			Id:   openapiutil.ToUUID(orgn.ID),
			Name: orgn.Name,
		})
	}

	return openapi.GetV1Orgns200JSONResponse{
		Organizations: orgns,
		NextCursor:    res.NextCursor,
	}, nil
}
