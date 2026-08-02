package finish_session

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/finish_session/dto"
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
	logger = logger.Set("pkg", "internal/transport/grpc/auth_service/rpc/finish_session")

	return &Handler{
		authorizer: authorizer,
		usecase:    usecase,
		logger:     logger,
	}
}

func (h *Handler) FinishSession(ctx context.Context, req *pb.FinishSessionRequest) (*emptypb.Empty, error) {
	auth, err := h.authorizer.Authorize(ctx)
	if err != nil {
		return nil, err
	}

	err = h.usecase.Do(ctx, &dto.Request{
		Auth:      auth,
		SessionID: req.SessionId,
	})
	if err != nil {
		return nil, grpcutil.ConvertError(err, h.logger)
	}

	return &emptypb.Empty{}, nil
}
