package register

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/register/dto"
	pb "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/generated/grpc/auth_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/grpcutil"
)

type Usecase interface {
	Do(ctx context.Context, req *dto.Request) error
}

type Handler struct {
	usecase Usecase
	logger  log.Logger
}

func NewHandler(
	usecase Usecase,
	logger log.Logger,
) *Handler {
	logger = logger.Set("pkg", "internal/transport/grpc/auth_service/rpc/register")

	return &Handler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *Handler) Register(ctx context.Context, req *pb.RegisterRequest) (*emptypb.Empty, error) {
	err := h.usecase.Do(ctx, &dto.Request{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, grpcutil.ConvertError(err, h.logger)
	}

	return &emptypb.Empty{}, nil
}
