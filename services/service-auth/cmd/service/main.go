package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/database/psql"
	psqlRepositories "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/database/psql/repositories"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/bcrypt"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/shutdown"
	httpEndpoints "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/http/endpoints"
)

var (
	logger     log.Logger
	shutdowner *shutdown.Shutdowner
	httpServer *http.Server
)

func main() {
	logger = log.NewLogger(nil).Set("pkg", "cmd/service")
	shutdowner = shutdown.NewManager(logger)
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
	psqlConn, err := psql.Connect(ctx, logger)
	if err != nil {
		return fmt.Errorf("connect psql: %w", err)
	}
	shutdowner.Add(psqlConn)

	hashService := bcrypt.NewService()

	// Repositories
	userRepository := psqlRepositories.NewUserRepository(psqlConn)
	emailVerificationRepository := psqlRepositories.NewEmailVerificationRepository(psqlConn)

	// Usecases
	registerUsecase := usecases.NewRegisterUsecase(hashService, userRepository, emailVerificationRepository, logger)
	verifyEmailUsecase := usecases.NewVerifyEmailUsecase(emailVerificationRepository, userRepository, logger)

	// HTTP Server
	httpRouter := http.NewServeMux()
	httpEndpoints.RegisterRegisterEndpoint(httpRouter, registerUsecase, logger)
	httpEndpoints.RegisterVerifyEmailEndpoint(httpRouter, verifyEmailUsecase, logger)

	httpServer = &http.Server{
		Addr:        ":8080",
		Handler:     httpRouter,
		BaseContext: func(net.Listener) context.Context { return ctx },
	}
	shutdowner.AddFunc(func() error {
		return httpServer.Shutdown(context.Background())
	})

	return nil
}
