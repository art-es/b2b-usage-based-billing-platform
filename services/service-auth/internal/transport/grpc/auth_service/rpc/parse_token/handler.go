package parse_token

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
	pb "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/generated/grpc/auth_service"
)

type authorizer interface {
	Authorize(ctx context.Context) (*jwt.Claims, error)
}

type Handler struct {
	authorizer authorizer
}

func NewHandler(authorizer authorizer) *Handler {
	return &Handler{
		authorizer: authorizer,
	}
}

func (h *Handler) ParseToken(ctx context.Context, _ *emptypb.Empty) (*pb.ParseTokenResponse, error) {
	claims, err := h.authorizer.Authorize(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.ParseTokenResponse{
		SessionId: claims.SessionID,
		UserId:    claims.UserID,
		OrgnId:    claims.OrgnID,
	}, nil
}
