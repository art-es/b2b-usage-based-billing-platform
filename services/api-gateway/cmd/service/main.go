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
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi"
)

var (
	logger     log.Logger
	shutdowner *shutdown.Shutdowner
	httpServer *http.Server
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
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Log(log.Error).
				Set("message", "http server listen error").
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
	)
	if err != nil {
		return fmt.Errorf("parse env vars: %w", err)
	}

	httpServer = &http.Server{
		Addr:        envs.Get(env.FieldApiGatewayAddr),
		Handler:     openapi.NewHandler(logger),
		BaseContext: func(net.Listener) context.Context { return ctx },
	}
	shutdowner.AddFunc(func() error {
		return httpServer.Shutdown(context.Background())
	})

	return nil
}
