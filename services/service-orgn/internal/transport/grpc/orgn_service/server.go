package orgn_service

import (
	"context"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/transport/grpc/orgn_service/rpc/update"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/generated/grpc/orgn_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/transport/grpc/orgn_service/rpc/create"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/transport/grpc/orgn_service/rpc/get"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/transport/grpc/orgn_service/rpc/get_by_id"
)

type getHandler interface {
	Get(context.Context, *pb.GetRequest) (*pb.GetResponse, error)
}

type getByIDHandler interface {
	GetById(context.Context, *pb.GetByIdRequest) (*pb.Orgn, error)
}

type createHandler interface {
	Create(context.Context, *pb.CreateRequest) (*pb.CreateResponse, error)
}

type updateHandler interface {
	Update(context.Context, *pb.UpdateRequest) (*emptypb.Empty, error)
}

type serverHandler struct {
	getHandler
	getByIDHandler
	createHandler
	updateHandler

	pb.UnsafeOrgnServiceServer
}

func NewServer(
	logger log.Logger,
	getUsecase get.Usecase,
	getByIDUsecase get_by_id.Usecase,
	createUsecase create.Usecase,
	updateUsecase update.Usecase,
) *grpc.Server {
	handler := &serverHandler{
		getHandler:     get.NewHandler(getUsecase, logger),
		getByIDHandler: get_by_id.NewHandler(getByIDUsecase, logger),
		createHandler:  create.NewHandler(createUsecase, logger),
		updateHandler:  update.NewHandler(updateUsecase, logger),
	}

	server := grpc.NewServer()
	pb.RegisterOrgnServiceServer(server, handler)
	return server
}
