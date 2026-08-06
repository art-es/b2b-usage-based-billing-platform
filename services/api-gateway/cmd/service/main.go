package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/data/env"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/pkg/shutdown"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/grpc/auth_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi"
)

const (
	envApiGatewayAddr  = "API_GATEWAY_ADDR"
	envAuthServiceAddr = "AUTH_SERVICE_ADDR"
)

var (
	logger        log.Logger
	shutdowner    *shutdown.Shutdowner
	runHTTPServer func() error
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
		Set("message", "service starting").
		Write()

	go func() {
		if err := runHTTPServer(); err != nil {
			logger.Log(log.Error).
				Set("message", "http server run error").
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
	err := env.CheckEmpty(envApiGatewayAddr, envAuthServiceAddr)
	if err != nil {
		return err
	}

	authService, err := auth_service.NewClient(os.Getenv(envAuthServiceAddr))
	if err != nil {
		return fmt.Errorf("create grpc client of auth service: %w", err)
	}
	shutdowner.Add(authService)

	openapiHandler := openapi.NewHandler(logger, authService)
	initHTTPServer(ctx, openapiHandler)

	return nil
}

func initHTTPServer(ctx context.Context, handler http.Handler) {
	server := &http.Server{
		Addr:        os.Getenv(envApiGatewayAddr),
		Handler:     handler,
		BaseContext: func(net.Listener) context.Context { return ctx },
	}

	shutdowner.AddFunc(func() error {
		return server.Shutdown(ctx)
	})

	runHTTPServer = func() error {
		return server.ListenAndServe()
	}
}
