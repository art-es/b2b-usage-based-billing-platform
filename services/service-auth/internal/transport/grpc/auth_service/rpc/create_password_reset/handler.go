package create_password_reset

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/generated/grpc/auth_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/grpcutil"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) CreatePasswordReset(context.Context, *pb.CreatePasswordResetRequest) (*emptypb.Empty, error) {
	return nil, grpcutil.Unimplemented()
}
