package auth_service

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
	pb "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/generated/grpc/auth_service"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/auth_service/rpc/change_password"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/auth_service/rpc/create_password_reset"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/auth_service/rpc/finish_all_sessions"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/auth_service/rpc/finish_session"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/auth_service/rpc/get_me"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/auth_service/rpc/get_sessions"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/auth_service/rpc/login"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/auth_service/rpc/refresh_session"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/auth_service/rpc/register"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/auth_service/rpc/resend_email_verification"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/auth_service/rpc/reset_password"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/auth_service/rpc/switch_orgn"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/auth_service/rpc/verify_email"
)

type (
	authorizer interface {
		Authorize(context.Context) (*jwt.Claims, error)
	}

	registerHandler interface {
		Register(context.Context, *pb.RegisterRequest) (*emptypb.Empty, error)
	}

	verifyEmailHandler interface {
		VerifyEmail(context.Context, *pb.VerifyEmailRequest) (*emptypb.Empty, error)
	}

	resendEmailVerificationHandler interface {
		ResendEmailVerification(context.Context, *pb.ResendEmailVerificationRequest) (*emptypb.Empty, error)
	}

	loginHandler interface {
		Login(context.Context, *pb.LoginRequest) (*pb.LoginResponse, error)
	}

	refreshSessionHandler interface {
		RefreshSession(context.Context, *pb.RefreshSessionRequest) (*pb.RefreshSessionResponse, error)
	}

	getSessionsHandler interface {
		GetSessions(context.Context, *pb.GetSessionsRequest) (*pb.GetSessionsResponse, error)
	}

	finishAllSessionsHandler interface {
		FinishAllSessions(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	}

	finishSessionHandler interface {
		FinishSession(context.Context, *pb.FinishSessionRequest) (*emptypb.Empty, error)
	}

	createPasswordResetHandler interface {
		CreatePasswordReset(context.Context, *pb.CreatePasswordResetRequest) (*emptypb.Empty, error)
	}

	resetPasswordHandler interface {
		ResetPassword(context.Context, *pb.ResetPasswordRequest) (*emptypb.Empty, error)
	}

	changePasswordHandler interface {
		ChangePassword(context.Context, *pb.ChangePasswordRequest) (*emptypb.Empty, error)
	}

	switchOrgnHandler interface {
		SwitchOrgn(context.Context, *pb.SwitchOrgnRequest) (*pb.SwitchOrgnResponse, error)
	}

	getMeHandler interface {
		GetMe(context.Context, *emptypb.Empty) (*pb.GetMeResponse, error)
	}
)

type serverHandler struct {
	registerHandler
	verifyEmailHandler
	resendEmailVerificationHandler
	loginHandler
	refreshSessionHandler
	getSessionsHandler
	finishAllSessionsHandler
	finishSessionHandler
	createPasswordResetHandler
	resetPasswordHandler
	changePasswordHandler
	switchOrgnHandler
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
	refreshSessionUsecase refresh_session.Usecase,
	getSessionsUsecase get_sessions.Usecase,
	finishAllSessionsUsecase finish_all_sessions.Usecase,
	finishSessionUsecase finish_session.Usecase,
	createPasswordResetUsecase create_password_reset.Usecase,
	resetPasswordUsecase reset_password.Usecase,
	changePasswordUsecase change_password.Usecase,
	switchOrgnUsecase switch_orgn.Usecase,
	getMeUsecase get_me.Usecase,
) *grpc.Server {
	handler := &serverHandler{
		registerHandler:                register.NewHandler(registerUsecase, logger),
		verifyEmailHandler:             verify_email.NewHandler(verifyEmailUsecase, logger),
		resendEmailVerificationHandler: resend_email_verification.NewHandler(resendEmailVerificationUsecase, logger),
		loginHandler:                   login.NewHandler(loginUsecase, logger),
		refreshSessionHandler:          refresh_session.NewHandler(refreshSessionUsecase, logger),
		getSessionsHandler:             get_sessions.NewHandler(authorizer, getSessionsUsecase, logger),
		finishAllSessionsHandler:       finish_all_sessions.NewHandler(authorizer, finishAllSessionsUsecase, logger),
		finishSessionHandler:           finish_session.NewHandler(authorizer, finishSessionUsecase, logger),
		createPasswordResetHandler:     create_password_reset.NewHandler(createPasswordResetUsecase, logger),
		resetPasswordHandler:           reset_password.NewHandler(resetPasswordUsecase, logger),
		changePasswordHandler:          change_password.NewHandler(authorizer, changePasswordUsecase, logger),
		switchOrgnHandler:              switch_orgn.NewHandler(authorizer, switchOrgnUsecase, logger),
		getMeHandler:                   get_me.NewHandler(authorizer, getMeUsecase, logger),
	}

	server := grpc.NewServer()
	pb.RegisterAuthServiceServer(server, handler)
	return server
}
