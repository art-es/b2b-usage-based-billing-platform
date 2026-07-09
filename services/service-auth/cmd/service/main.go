package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/database/psql"
	psqlRepositories "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/database/psql/repositories"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/bcrypt"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/cmdutil"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	httpEndpoints "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/http/endpoints"
)

var (
	logger     log.Logger
	gsManager  *cmdutil.GSManager
	httpServer *http.Server
)

func main() {
	logger = log.NewLogger(nil).Set("pkg", "cmd/service")
	gsManager = cmdutil.NewGSManager(logger)

	if err := build(); err != nil {
		logger.Log(log.Error).
			Set("message", "service build error").
			Write()
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)

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

	gsManager.Shutdown()
}

func build() error {
	psqlConn, err := psql.Connect(logger)
	if err != nil {
		return fmt.Errorf("connect psql: %w", err)
	}
	gsManager.Add(psqlConn)

	hashService := bcrypt.NewService()

	// Repositories
	userRepository := psqlRepositories.NewUserRepository(psqlConn)

	// Usecases
	registerUsecase := usecases.NewRegisterUsecase(userRepository, hashService)

	// HTTP Endpoints
	httpRouter := http.NewServeMux()
	httpEndpoints.RegisterRegisterEndpoint(httpRouter, registerUsecase, logger)

	httpServer = &http.Server{Addr: ":8080", Handler: httpRouter}
	gsManager.AddFunc(func() error {
		return httpServer.Shutdown(context.Background())
	})

	return nil
}
