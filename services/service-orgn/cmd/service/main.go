package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/usecases"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/data/env"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/data/psql"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/data/psql/repositories"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/pkg/shutdown"
	grpc_orgn_service "github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/transport/grpc/orgn_service"
	"google.golang.org/grpc"
)

const (
	envPSQLAddr         = "PSQL_ADDR"
	envOrgnServiceAddr  = "ORGN_SERVICE_ADDR"
	envOrgnCursorSecret = "ORGN_CURSOR_SECRET"
)

var (
	logger        log.Logger
	shutdowner    *shutdown.Shutdowner
	runGRPCServer func() error
)

func main() {
	logger = log.NewLogger(nil).Set("pkg", "cmd/service")

	shutdowner = shutdown.New(logger)
	defer shutdowner.Shutdown()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	if err := build(ctx); err != nil {
		logger.Log(log.Error).
			Set("message", "build error").
			Write()
		return
	}

	logger.Log(log.Info).
		Set("message", "service started").
		Write()

	go func() {
		if err := runGRPCServer(); err != nil {
			logger.Log(log.Error).
				Set("message", "grpc server run error").
				Set("error", err.Error()).
				Write()
		}
		stop()
	}()

	<-ctx.Done()

	logger.Log(log.Info).
		Set("message", "service finished").
		Write()
}

func build(ctx context.Context) error {
	err := env.CheckEmpty(envPSQLAddr, envOrgnCursorSecret)
	if err != nil {
		return err
	}

	// Clients
	psqlConn, err := psql.Connect(ctx, os.Getenv(envPSQLAddr), logger)
	if err != nil {
		return fmt.Errorf("connect psql: %w", err)
	}
	shutdowner.Add(psqlConn)

	// Repositories
	orgnRepository := repositories.NewOrgnRepository(psqlConn)

	// Usecases
	getUsecase := usecases.NewGetUsecase(orgnRepository, os.Getenv(envOrgnCursorSecret), logger)
	getByIDUsecase := usecases.NewGetByIDUsecase(orgnRepository)
	createUsecase := usecases.NewCreateUsecase(orgnRepository)

	// GRPC Server
	grpcServer := grpc_orgn_service.NewServer(
		logger,
		getUsecase,
		getByIDUsecase,
		createUsecase,
	)
	initGRPCServer(grpcServer)

	return nil
}

func initGRPCServer(server *grpc.Server) error {
	addr := os.Getenv(envOrgnServiceAddr)
	if addr == "" {
		addr = ":8080"
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen grpc server port: %w", err)
	}

	shutdowner.Add(listener)
	shutdowner.AddFunc(func() error {
		server.GracefulStop()
		return nil
	})

	runGRPCServer = func() error {
		return server.Serve(listener)
	}

	return nil
}
