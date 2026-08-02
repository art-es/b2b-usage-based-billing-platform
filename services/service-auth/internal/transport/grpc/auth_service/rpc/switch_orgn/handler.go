package switch_orgn

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/switch_orgn/dto"
	pb "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/generated/grpc/auth_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/grpcutil"
)

type authorizer interface {
	Authorize(ctx context.Context) (*jwt.Claims, error)
}

type Usecase interface {
	Do(ctx context.Context, req *dto.Request) error
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
	logger = logger.Set("pkg", "internal/transport/grpc/auth_service/rpc/switch_orgn")

	return &Handler{
		authorizer: authorizer,
		usecase:    usecase,
		logger:     logger,
	}
}

func (h *Handler) SwitchOrgn(ctx context.Context, req *pb.SwitchOrgnRequest) (*emptypb.Empty, error) {
	auth, err := h.authorizer.Authorize(ctx)
	if err != nil {
		return nil, err
	}

	err = h.usecase.Do(ctx, &dto.Request{
		Auth:   auth,
		OrgnID: req.OrgnId,
	})
	if err != nil {
		return nil, grpcutil.ConvertError(err, h.logger)
	}

	return &emptypb.Empty{}, nil
}
