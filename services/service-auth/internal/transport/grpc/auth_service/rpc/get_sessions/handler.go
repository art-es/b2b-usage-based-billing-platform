package get_sessions

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/get_sessions/dto"
	pb "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/generated/grpc/auth_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/grpcutil"
)

type authorizer interface {
	Authorize(ctx context.Context) (*jwt.Claims, error)
}

type Usecase interface {
	Do(ctx context.Context, req *dto.Request) (*dto.Response, error)
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
	logger = logger.Set("pkg", "internal/transport/grpc/auth_service/rpc/get_sessions")

	return &Handler{
		authorizer: authorizer,
		usecase:    usecase,
		logger:     logger,
	}
}

func (h *Handler) GetSessions(ctx context.Context, req *pb.GetSessionsRequest) (*pb.GetSessionsResponse, error) {
	claims, err := h.authorizer.Authorize(ctx)
	if err != nil {
		return nil, err
	}

	ucRes, err := h.usecase.Do(ctx, &dto.Request{
		Auth:   claims,
		Cursor: req.Cursor,
	})
	if err != nil {
		return nil, grpcutil.ConvertError(err, h.logger)
	}

	sessions := make([]*pb.Session, 0, len(ucRes.Sessions))
	for _, session := range ucRes.Sessions {
		sessions = append(sessions, convertSession(session))
	}

	return &pb.GetSessionsResponse{
		Sessions:   sessions,
		NextCursor: ucRes.NextCursor,
	}, nil
}

func convertSession(in *dto.Session) *pb.Session {
	return &pb.Session{
		Id:        in.ID,
		CreatedAt: timestamppb.New(in.CreatedAt),
	}
}
