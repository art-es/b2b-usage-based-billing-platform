package main

import (
	"context"
	"fmt"
	"net"
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
	grpc_auth_service "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/auth_service"
	grpc_auth "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/auth_service/auth"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/orgn_service"
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
		env.Required(env.FieldOrgnServiceAddr),
		env.Required(env.FieldJwtSecret),
		env.Required(env.FieldRefreshTokenSecret),
		env.Required(env.FieldSessionsCursorSecret),
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

	orgnService, err := orgn_service.NewClient(envs.Get(env.FieldOrgnServiceAddr))
	if err != nil {
		return fmt.Errorf("connect orgn-service: %w", err)
	}
	shutdowner.Add(orgnService)

	// Utils
	timeService := time.NewService()
	uuidService := uuid.NewService()
	jwtService := jwt.NewService(timeService, logger)
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
		userRepository, envs.Get(env.FieldJwtSecret), envs.Get(env.FieldRefreshTokenSecret), logger,
	)
	refreshSessionUsecase := usecases.NewRefreshSessionUsecase(
		jwtService, hmacSha256Service, timeService, uuidService, sessionRepository,
		envs.Get(env.FieldJwtSecret), envs.Get(env.FieldRefreshTokenSecret), logger,
	)
	getSessionsUsecase := usecases.NewGetSessionsUsecase(sessionRepository, hmacSha256Service, envs.Get(env.FieldSessionsCursorSecret))
	finishAllSessionsUsecase := usecases.NewFinishAllSessionsUsecase(sessionRepository)
	finishSessionUsecase := usecases.NewFinishSessionUsecase(sessionRepository)
	switchOrgnUsecase := usecases.NewSwitchOrgnUsecase(sessionRepository, orgnService)
	getMeUsecase := usecases.NewGetMeUsecase(userRepository, orgnService)

	// GRPC Server
	grpcServerListener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return fmt.Errorf("listen grpc server port: %w", err)
	}
	shutdowner.Add(grpcServerListener)

	grpcAuthorizer := grpc_auth.NewAuthorizer(jwtService, envs.Get(env.FieldJwtSecret), logger)
	grpcServer := grpc_auth_service.NewServer(
		logger,
		grpcAuthorizer,
		registerUsecase,
		verifyEmailUsecase,
		resendEmailVerificationsUsecase,
		loginUsecase,
		refreshSessionUsecase,
		getSessionsUsecase,
		finishAllSessionsUsecase,
		finishSessionUsecase,
		switchOrgnUsecase,
		getMeUsecase,
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
