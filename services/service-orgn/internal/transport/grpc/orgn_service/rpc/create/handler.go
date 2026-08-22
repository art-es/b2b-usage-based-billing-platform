package create

import (
	"context"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/usecases/create/dto"
	pb "github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/generated/grpc/orgn_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/transport/grpc/grpcutil"
)

type Usecase interface {
	Do(ctx context.Context, req *dto.Request) (*dto.Response, error)
}

type Handler struct {
	usecase Usecase
	logger  log.Logger
}

func NewHandler(
	usecase Usecase,
	logger log.Logger,
) *Handler {
	logger = logger.Set("pkg", "internal/transport/grpc/orgn_service/rpc/create")

	return &Handler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *Handler) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	res, err := h.usecase.Do(ctx, &dto.Request{
		UserID: req.UserId,
		Name:   req.Name,
	})
	if err != nil {
		return nil, grpcutil.ConvertError(err, h.logger)
	}

	return &pb.CreateResponse{
		OrgnId: res.OrgnID,
	}, nil
}
