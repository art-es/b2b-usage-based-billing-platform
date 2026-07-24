package auth_service

import (
	"context"

	"google.golang.org/grpc"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
	pb "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/generated/grpc/auth_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/auth_service/rpc/get_me"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/auth_service/rpc/login"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/auth_service/rpc/register"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/auth_service/rpc/resend_email_verification"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/auth_service/rpc/verify_email"
)

type (
	authorizer interface {
		Authorize(ctx context.Context) (*jwt.Claims, error)
	}

	registerHandler interface {
		Register(context.Context, *pb.RegisterRequest) (*pb.Empty, error)
	}

	verifyEmailHandler interface {
		VerifyEmail(context.Context, *pb.VerifyEmailRequest) (*pb.Empty, error)
	}

	resendEmailVerificationHandler interface {
		ResendEmailVerification(context.Context, *pb.ResendEmailVerificationRequest) (*pb.Empty, error)
	}

	loginHandler interface {
		Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
	}

	getMeHandler interface {
		GetMe(context.Context, *pb.Empty) (*pb.GetMeResponse, error)
	}
)

type serverHandler struct {
	registerHandler
	verifyEmailHandler
	resendEmailVerificationHandler
	loginHandler
	getMeHandler
	pb.UnsafeAuthServiceServer
}

func NewServer(
	logger log.Logger,
	authorizer authorizer,
	registerUsecase register.Usecase,
	verifyEmailUsecase verify_email.Usecase,
	resendEmailVerificationUsecase resend_email_verification.Usecase,
	loginUsecase login.Usecase,
	getMeUsecase get_me.Usecase,
) *grpc.Server {
	handler := &serverHandler{
		registerHandler:                register.NewHandler(registerUsecase, logger),
		verifyEmailHandler:             verify_email.NewHandler(verifyEmailUsecase, logger),
		resendEmailVerificationHandler: resend_email_verification.NewHandler(resendEmailVerificationUsecase, logger),
		loginHandler:                   login.NewHandler(loginUsecase, logger),
		getMeHandler:                   get_me.NewHandler(authorizer, getMeUsecase, logger),
	}

	server := grpc.NewServer()
	pb.RegisterAuthServiceServer(server, handler)
	return server
}
