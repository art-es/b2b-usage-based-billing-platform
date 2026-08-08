package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

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

const (
	envPSQLAddr             = "PSQL_ADDR"
	envAuthServiceAddr      = "AUTH_SERVICE_ADDR"
	envOrgnServiceAddr      = "ORGN_SERVICE_ADDR"
	envJWTSecret            = "JWT_SECRET"
	envRefreshTokenSecret   = "REFRESH_TOKEN_SECRET"
	envSessionsCursorSecret = "SESSIONS_CURSOR_SECRET"
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
	err := env.CheckEmpty(envPSQLAddr, envOrgnServiceAddr, envJWTSecret, envRefreshTokenSecret, envSessionsCursorSecret)
	if err != nil {
		return err
	}

	// Clients
	psqlConn, err := psql.Connect(ctx, os.Getenv(envPSQLAddr), logger)
	if err != nil {
		return fmt.Errorf("connect psql: %w", err)
	}
	shutdowner.Add(psqlConn)

	orgnService, err := orgn_service.NewClient(os.Getenv(envOrgnServiceAddr))
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
	passwordResetRepository := repositories.NewPasswordResetRepository(psqlConn)

	// Usecases
	registerUsecase := usecases.NewRegisterUsecase(passwordHashService, userRepository, emailVerificationRepository, logger)
	verifyEmailUsecase := usecases.NewVerifyEmailUsecase(emailVerificationRepository, userRepository, logger)
	resendEmailVerificationsUsecase := usecases.NewResendEmailVerificationUsecase(emailVerificationRepository)
	loginUsecase := usecases.NewLoginUsecase(
		jwtService, hmacSha256Service, passwordHashService, timeService, uuidService, sessionRepository,
		userRepository, os.Getenv(envJWTSecret), os.Getenv(envRefreshTokenSecret), logger,
	)
	refreshSessionUsecase := usecases.NewRefreshSessionUsecase(
		jwtService, hmacSha256Service, timeService, uuidService, sessionRepository,
		os.Getenv(envJWTSecret), os.Getenv(envRefreshTokenSecret), logger,
	)
	getSessionsUsecase := usecases.NewGetSessionsUsecase(sessionRepository, hmacSha256Service, os.Getenv(envSessionsCursorSecret))
	finishAllSessionsUsecase := usecases.NewFinishAllSessionsUsecase(sessionRepository)
	finishSessionUsecase := usecases.NewFinishSessionUsecase(sessionRepository)
	createPasswordResetUsecase := usecases.NewCreatePasswordResetUsecase(userRepository, passwordResetRepository)
	resetPasswordUsecase := usecases.NewResetPasswordUsecase(passwordResetRepository, userRepository, passwordHashService, logger)
	changePasswordUsecase := usecases.NewChangePasswordUsecase(userRepository, passwordHashService)
	switchOrgnUsecase := usecases.NewSwitchOrgnUsecase(sessionRepository, orgnService, jwtService, os.Getenv(envJWTSecret), logger)
	getMeUsecase := usecases.NewGetMeUsecase(userRepository, orgnService)

	// GRPC Server
	grpcAuthorizer := grpc_auth.NewAuthorizer(jwtService, os.Getenv(envJWTSecret), logger)
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
		createPasswordResetUsecase,
		resetPasswordUsecase,
		changePasswordUsecase,
		switchOrgnUsecase,
		getMeUsecase,
	)
	initGRPCServer(grpcServer)

	return nil
}

func initGRPCServer(server *grpc.Server) error {
	addr := os.Getenv(envAuthServiceAddr)
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
		server.Serve(listener)
		return nil
	}

	return nil
}
