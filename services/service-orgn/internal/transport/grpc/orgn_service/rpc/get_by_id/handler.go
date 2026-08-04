package get_by_id

import (
	"context"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/domains/orgn"
	pb "github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/generated/grpc/orgn_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/transport/grpc/grpcutil"
)

type Usecase interface {
	Do(ctx context.Context, orgnID string) (*orgn.Orgn, error)
}

type Handler struct {
	usecase Usecase
	logger  log.Logger
}

func NewHandler(
	usecase Usecase,
	logger log.Logger,
) *Handler {
	logger = logger.Set("pkg", "internal/transport/grpc/orgn_service/rpc/get_by_id")

	return &Handler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *Handler) GetById(ctx context.Context, req *pb.GetByIdRequest) (*pb.Orgn, error) {
	org, err := h.usecase.Do(ctx, req.Id)
	if err != nil {
		return nil, grpcutil.ConvertError(err, h.logger)
	}

	return &pb.Orgn{
		Id:     org.ID,
		Name:   org.Name,
		UserId: org.UserID,
	}, nil
}
