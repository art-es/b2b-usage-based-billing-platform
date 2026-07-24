package get_me

import (
	"context"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/get_me/dto"
	pb "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/generated/grpc/auth_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/grpcutil"
)

type authorizer interface {
	Authorize(ctx context.Context) (*jwt.Claims, error)
}

type Usecase interface {
	Do(ctx context.Context, claims *jwt.Claims) (*dto.Response, error)
}

type Handler struct {
	authorizer authorizer
	usecase    Usecase
	logger     log.Logger
}

func NewHandler(
	authorizer authorizer,
	usecase Usecase,
	logger log.Logger,
) *Handler {
	logger = logger.Set("pkg", "internal/transport/http/endpoints/get_me")

	return &Handler{
		authorizer: authorizer,
		usecase:    usecase,
		logger:     logger,
	}
}

func (h *Handler) GetMe(ctx context.Context, _ *pb.Empty) (*pb.GetMeResponse, error) {
	token, err := h.authorizer.Authorize(ctx)
	if err != nil {
		return nil, err
	}

	ucRes, err := h.usecase.Do(ctx, token)
	if err != nil {
		return nil, grpcutil.ConvertError(err, h.logger)
	}

	pbRes := &pb.GetMeResponse{
		SessionId: ucRes.SessionID,
		Name:      ucRes.User.Name,
		Email:     ucRes.User.Email,
	}

	if ucRes.Orgn != nil {
		pbRes.Orgn = &pb.GetMeResponseOrgn{
			Id:   ucRes.Orgn.ID,
			Name: ucRes.Orgn.Name,
		}
	}

	return pbRes, nil
}
