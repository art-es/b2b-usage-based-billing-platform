package get

import (
	"context"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/domains/orgn"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/usecases/get/dto"
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
	logger = logger.Set("pkg", "internal/transport/grpc/orgn_service/rpc/get")

	return &Handler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *Handler) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	res, err := h.usecase.Do(ctx, &dto.Request{
		UserID: req.UserId,
		Cursor: req.Cursor,
	})
	if err != nil {
		return nil, grpcutil.ConvertError(err, h.logger)
	}

	return &pb.GetResponse{
		Orgns:      convertOrgns(res.List),
		NextCursor: res.NextCursor,
	}, nil
}

func convertOrgns(in []*orgn.Orgn) []*pb.Orgn {
	out := make([]*pb.Orgn, 0, len(in))
	for _, o := range in {
		out = append(out, &pb.Orgn{
			Id:     o.ID,
			Name:   o.Name,
			UserId: o.UserID,
		})
	}
	return out
}
