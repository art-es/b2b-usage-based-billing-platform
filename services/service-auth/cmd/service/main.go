package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/data/env"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/data/psql"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/data/psql/repositories"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/bcrypt"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/hmac_sha256"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/jwt"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/shutdown"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/time"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/uuid"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/http/endpoints"
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
	envs, err := env.ParseVars(
		env.Required(env.FieldPsqlUrl),
		env.Required(env.FieldJwtSecret),
		env.Required(env.FieldRefreshTokenHashSecret),
	)
	if err != nil {
		return fmt.Errorf("parse env vars: %w", err)
	}

	psqlConn, err := psql.Connect(ctx, envs.Get(env.FieldPsqlUrl), logger)
	if err != nil {
		return fmt.Errorf("connect psql: %w", err)
	}
	shutdowner.Add(psqlConn)

	timeService := time.NewService()
	uuidService := uuid.NewService()
	jwtService := jwt.NewService(logger)
	passwordHashService := bcrypt.NewService()
	hmacSha256Service := hmac_sha256.NewService()

	// Repositories
	userRepository := repositories.NewUserRepository(psqlConn)
	emailVerificationRepository := repositories.NewEmailVerificationRepository(psqlConn)
	sessionRepository := repositories.NewSessionsRepository(psqlConn)

	// Usecases
	registerUsecase := usecases.NewRegisterUsecase(passwordHashService, userRepository, emailVerificationRepository, logger)
	verifyEmailUsecase := usecases.NewVerifyEmailUsecase(emailVerificationRepository, userRepository, logger)
	resendEmailVerificationsUsecase := usecases.NewResendEmailVerificationUsecase(emailVerificationRepository)
	loginUsecase := usecases.NewLoginUsecase(
		jwtService, hmacSha256Service, passwordHashService, timeService, uuidService, sessionRepository,
		userRepository, envs.Get(env.FieldJwtSecret), envs.Get(env.FieldRefreshTokenHashSecret), logger,
	)

	// HTTP Server
	httpRouter := http.NewServeMux()
	endpoints.BindRegister(httpRouter, registerUsecase, logger)
	endpoints.BindVerifyEmail(httpRouter, verifyEmailUsecase, logger)
	endpoints.BindResendEmailVerification(httpRouter, resendEmailVerificationsUsecase, logger)
	endpoints.BindLogin(httpRouter, loginUsecase, logger)

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
