package openapi

import (
	"context"
	"net/http"

	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/auth"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/delete_v1_auth_sessions"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/delete_v1_auth_sessions_session_id"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/get_v1_auth_sessions"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/get_v1_me"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/post_v1_auth_email_resend_verification"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/post_v1_auth_email_verify"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/post_v1_auth_login"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/post_v1_auth_refresh"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/post_v1_auth_register"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/post_v1_auth_switch_orgn"
)

type AuthService interface {
	post_v1_auth_register.AuthService
	post_v1_auth_email_verify.AuthService
	post_v1_auth_email_resend_verification.AuthService
	post_v1_auth_login.AuthService
	post_v1_auth_refresh.AuthService
	get_v1_auth_sessions.AuthService
	delete_v1_auth_sessions.AuthService
	delete_v1_auth_sessions_session_id.AuthService
	post_v1_auth_switch_orgn.AuthService
	get_v1_me.AuthService
}

// Endpoint handlers
type (
	postV1AuthRegisterHandler interface {
		PostV1AuthRegister(context.Context, openapi.PostV1AuthRegisterRequestObject) (openapi.PostV1AuthRegisterResponseObject, error)
	}

	postV1AuthEmailVerifyHandler interface {
		PostV1AuthEmailVerify(context.Context, openapi.PostV1AuthEmailVerifyRequestObject) (openapi.PostV1AuthEmailVerifyResponseObject, error)
	}

	postV1AuthEmailResendVerificationHandler interface {
		PostV1AuthEmailResendVerification(context.Context, openapi.PostV1AuthEmailResendVerificationRequestObject) (openapi.PostV1AuthEmailResendVerificationResponseObject, error)
	}

	postV1AuthLoginHandler interface {
		PostV1AuthLogin(context.Context, openapi.PostV1AuthLoginRequestObject) (openapi.PostV1AuthLoginResponseObject, error)
	}

	postV1AuthRefreshHandler interface {
		PostV1AuthRefresh(context.Context, openapi.PostV1AuthRefreshRequestObject) (openapi.PostV1AuthRefreshResponseObject, error)
	}

	getV1AuthSessionsHandler interface {
		GetV1AuthSessions(context.Context, openapi.GetV1AuthSessionsRequestObject) (openapi.GetV1AuthSessionsResponseObject, error)
	}

	deleteV1AuthSessionsHandler interface {
		DeleteV1AuthSessions(context.Context, openapi.DeleteV1AuthSessionsRequestObject) (openapi.DeleteV1AuthSessionsResponseObject, error)
	}

	deleteV1AuthSessionsSessionIdHandler interface {
		DeleteV1AuthSessionsSessionId(ctx context.Context, req openapi.DeleteV1AuthSessionsSessionIdRequestObject) (openapi.DeleteV1AuthSessionsSessionIdResponseObject, error)
	}

	postV1AuthSwitchOrgnHandler interface {
		PostV1AuthSwitchOrgn(context.Context, openapi.PostV1AuthSwitchOrgnRequestObject) (openapi.PostV1AuthSwitchOrgnResponseObject, error)
	}

	getV1MeHandler interface {
		GetV1Me(context.Context, openapi.GetV1MeRequestObject) (openapi.GetV1MeResponseObject, error)
	}
)

type serverHandler struct {
	postV1AuthRegisterHandler
	postV1AuthEmailVerifyHandler
	postV1AuthEmailResendVerificationHandler
	postV1AuthLoginHandler
	postV1AuthRefreshHandler
	getV1AuthSessionsHandler
	deleteV1AuthSessionsHandler
	deleteV1AuthSessionsSessionIdHandler
	postV1AuthSwitchOrgnHandler
	getV1MeHandler
}

func NewHandler(
	logger log.Logger,
	authService AuthService,
) http.Handler {
	logger = logger.Set("pkg", "internal/transport/http/openapi")

	hand := &serverHandler{
		postV1AuthRegisterHandler:                post_v1_auth_register.NewHandler(authService),
		postV1AuthEmailVerifyHandler:             post_v1_auth_email_verify.NewHandler(authService),
		postV1AuthEmailResendVerificationHandler: post_v1_auth_email_resend_verification.NewHandler(authService),
		postV1AuthLoginHandler:                   post_v1_auth_login.NewHandler(authService),
		postV1AuthRefreshHandler:                 post_v1_auth_refresh.NewHandler(authService),
		getV1AuthSessionsHandler:                 get_v1_auth_sessions.NewHandler(authService),
		deleteV1AuthSessionsHandler:              delete_v1_auth_sessions.NewHandler(authService),
		deleteV1AuthSessionsSessionIdHandler:     delete_v1_auth_sessions_session_id.NewHandler(authService),
		postV1AuthSwitchOrgnHandler:              post_v1_auth_switch_orgn.NewHandler(authService),
		getV1MeHandler:                           get_v1_me.NewHandler(authService),
	}

	opts := openapi.StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  getRequestErrorHandlerFunc(),
		ResponseErrorHandlerFunc: getResponseErrorHandlerFunc(logger),
	}

	mdlw := []openapi.StrictMiddlewareFunc{
		parseAccessToken,
	}

	return openapi.Handler(openapi.NewStrictHandlerWithOptions(hand, mdlw, opts))
}

func parseAccessToken(f openapi.StrictHandlerFunc, _ string) openapi.StrictHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (any, error) {
		if val := r.Header.Get("Authorization"); val != "" {
			ctx = auth.Set(ctx, val)
		}

		return f(ctx, w, r, request)
	}
}
