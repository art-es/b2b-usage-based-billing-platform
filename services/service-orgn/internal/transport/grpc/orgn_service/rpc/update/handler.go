package update

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/usecases/update/dto"
	pb "github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/generated/grpc/orgn_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/transport/grpc/grpcutil"
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
	logger = logger.Set("pkg", "internal/transport/grpc/orgn_service/rpc/update")

	return &Handler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *Handler) Update(ctx context.Context, req *pb.UpdateRequest) (*emptypb.Empty, error) {
	err := h.usecase.Do(ctx, &dto.Request{
		UserID: req.UserId,
		OrgnID: req.OrgnId,
		Name:   req.Name,
	})
	if err != nil {
		return nil, grpcutil.ConvertError(err, h.logger)
	}

	return &emptypb.Empty{}, nil
}
