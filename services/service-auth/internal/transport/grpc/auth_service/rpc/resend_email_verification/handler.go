package resend_email_verification

import (
	"context"

	pb "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/generated/grpc/auth_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/grpcutil"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Usecase interface {
	Do(ctx context.Context, email string) error
}

type Handler struct {
	usecase Usecase
	logger  log.Logger
}

func NewHandler(
	usecase Usecase,
	logger log.Logger,
) *Handler {
	logger = logger.Set("pkg", "internal/transport/grpc/auth_service/rpc/resend_email_verification")

	return &Handler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *Handler) ResendEmailVerification(ctx context.Context, req *pb.ResendEmailVerificationRequest) (*emptypb.Empty, error) {
	err := h.usecase.Do(ctx, req.Email)
	if err != nil {
		return nil, grpcutil.ConvertError(err, h.logger)
	}

	return &emptypb.Empty{}, nil
}
