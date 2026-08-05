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
	envs, err := env.ParseVars(
		env.Required(env.FieldApiGatewayAddr),
		env.Required(env.FieldAuthServiceAddr),
	)
	if err != nil {
		return fmt.Errorf("parse env vars: %w", err)
	}

	authService, err := auth_service.NewClient(envs.Get(env.FieldAuthServiceAddr))
	if err != nil {
		return fmt.Errorf("create grpc client of auth service: %w", err)
	}
	shutdowner.Add(authService)

	httpServer := &http.Server{
		Addr: envs.Get(env.FieldApiGatewayAddr),
		Handler: openapi.NewHandler(
			logger,
			authService,
		),
		BaseContext: func(net.Listener) context.Context { return ctx },
	}
	shutdowner.AddFunc(func() error {
		return httpServer.Shutdown(context.Background())
	})

	runHTTPServer = func() error {
		return httpServer.ListenAndServe()
	}

	return nil
}
