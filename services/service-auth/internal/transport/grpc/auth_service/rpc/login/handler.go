package login

import (
	"context"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/login/dto"
	pb "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/generated/grpc/auth_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/grpcutil"
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
	logger = logger.Set("pkg", "internal/transport/grpc/auth_service/rpc/login")

	return &Handler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *Handler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	res, err := h.usecase.Do(ctx, &dto.Request{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, grpcutil.ConvertError(err, h.logger)
	}

	return &pb.LoginResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}
