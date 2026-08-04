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
	envs, err := env.ParseVars(
		env.Required(env.FieldPsqlAddr),
		env.FieldOrgnServiceAddr,
	)
	if err != nil {
		return fmt.Errorf("parse env vars: %w", err)
	}

	// Clients
	psqlConn, err := psql.Connect(ctx, envs.Get(env.FieldPsqlAddr), logger)
	if err != nil {
		return fmt.Errorf("connect psql: %w", err)
	}
	shutdowner.Add(psqlConn)

	// Repositories
	orgnRepository := repositories.NewOrgnRepository(psqlConn)

	// Usecases
	getByIDUsecase := usecases.NewGetByIDUsecase(orgnRepository)

	// GRPC Server
	grpcServerAddr := envs.Get(env.FieldOrgnServiceAddr)
	if grpcServerAddr == "" {
		grpcServerAddr = ":8080"
	}

	grpcServerListener, err := net.Listen("tcp", grpcServerAddr)
	if err != nil {
		return fmt.Errorf("listen grpc server port: %w", err)
	}
	shutdowner.Add(grpcServerListener)

	grpcServer := grpc_orgn_service.NewServer(
		logger,
		getByIDUsecase,
	)
	shutdowner.AddFunc(func() error {
		grpcServer.GracefulStop()
		return nil
	})

	runGRPCServer = func() error {
		return grpcServer.Serve(grpcServerListener)
	}

	return nil
}
