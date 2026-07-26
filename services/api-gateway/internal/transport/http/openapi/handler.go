package openapi

import (
	"context"
	"net/http"

	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/post_v1_auth_email_verify"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/post_v1_auth_login"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/post_v1_auth_register"
)

type authService interface {
	post_v1_auth_register.AuthService
	post_v1_auth_email_verify.AuthService
	post_v1_auth_login.AuthService
}

// Endpoint handlers
type (
	postV1AuthRegisterHandler interface {
		PostV1AuthRegister(context.Context, openapi.PostV1AuthRegisterRequestObject) (openapi.PostV1AuthRegisterResponseObject, error)
	}

	postV1AuthEmailVerify interface {
		PostV1AuthEmailVerify(context.Context, openapi.PostV1AuthEmailVerifyRequestObject) (openapi.PostV1AuthEmailVerifyResponseObject, error)
	}

	postV1AuthLoginHandler interface {
		PostV1AuthLogin(context.Context, openapi.PostV1AuthLoginRequestObject) (openapi.PostV1AuthLoginResponseObject, error)
	}
)

type serverHandler struct {
	postV1AuthRegisterHandler
	postV1AuthEmailVerify
	postV1AuthLoginHandler
}

func NewHandler(
	logger log.Logger,
	authService authService,
) http.Handler {
	logger = logger.Set("pkg", "internal/transport/http/openapi")

	hand := &serverHandler{
		postV1AuthRegisterHandler: post_v1_auth_register.NewHandler(authService),
		postV1AuthEmailVerify:     post_v1_auth_email_verify.NewHandler(authService),
		postV1AuthLoginHandler:    post_v1_auth_login.NewHandler(authService),
	}

	opts := openapi.StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  getRequestErrorHandlerFunc(),
		ResponseErrorHandlerFunc: getResponseErrorHandlerFunc(logger),
	}

	return openapi.Handler(openapi.NewStrictHandlerWithOptions(hand, nil, opts))
}
