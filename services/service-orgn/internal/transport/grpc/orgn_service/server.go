package orgn_service

import (
	"context"

	"google.golang.org/grpc"

	pb "github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/generated/grpc/orgn_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/transport/grpc/orgn_service/rpc/get_by_id"
)

type getByIDHandler interface {
	GetById(context.Context, *pb.GetByIdRequest) (*pb.Orgn, error)
}

type serverHandler struct {
	getByIDHandler

	pb.UnsafeOrgnServiceServer
}

func NewServer(
	logger log.Logger,
	getByIDUsecase get_by_id.Usecase,
) *grpc.Server {
	handler := &serverHandler{
		getByIDHandler: get_by_id.NewHandler(getByIDUsecase, logger),
	}

	server := grpc.NewServer()
	pb.RegisterOrgnServiceServer(server, handler)
	return server
}
