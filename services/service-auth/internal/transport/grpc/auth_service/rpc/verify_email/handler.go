package verify_email

import (
	"context"

	pb "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/generated/grpc/auth_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/grpcutil"
)

type Usecase interface {
	Do(ctx context.Context, token string) error
}

type Handler struct {
	usecase Usecase
	logger  log.Logger
}

func NewHandler(
	usecase Usecase,
	logger log.Logger,
) *Handler {
	logger = logger.Set("pkg", "internal/transport/http/endpoints/verify_email")

	return &Handler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *Handler) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.Empty, error) {
	err := h.usecase.Do(ctx, req.Token)
	if err != nil {
		return nil, grpcutil.ConvertError(err, h.logger)
	}

	return &pb.Empty{}, nil
}
