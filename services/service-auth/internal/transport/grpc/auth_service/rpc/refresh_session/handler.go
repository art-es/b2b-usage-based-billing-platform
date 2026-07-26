package refresh_session

import (
	"context"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/refresh_session/dto"
	pb "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/generated/grpc/auth_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/grpcutil"
)

type Usecase interface {
	Do(ctx context.Context, token string) (*dto.Response, error)
}

type Handler struct {
	usecase Usecase
	logger  log.Logger
}

func NewHandler(
	usecase Usecase,
	logger log.Logger,
) *Handler {
	logger = logger.Set("pkg", "internal/transport/grpc/auth_service/rpc/refresh_session")

	return &Handler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *Handler) RefreshSession(ctx context.Context, req *pb.RefreshSessionRequest) (*pb.RefreshSessionResponse, error) {
	res, err := h.usecase.Do(ctx, req.Token)
	if err != nil {
		return nil, grpcutil.ConvertError(err, h.logger)
	}

	return &pb.RefreshSessionResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}
