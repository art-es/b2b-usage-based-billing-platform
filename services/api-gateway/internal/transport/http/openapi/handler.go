package openapi

import (
	"context"
	"net/http"

	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/auth"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/delete_v1_auth_sessions"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/delete_v1_auth_sessions_session_id"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/delete_v1_orgns_orgn_id"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/delete_v1_orgns_orgn_id_api_keys_key_id"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/delete_v1_orgns_orgn_id_members_member_id"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/delete_v1_orgns_orgn_id_tariffs_tariff_id"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/get_v1_auth_sessions"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/get_v1_billing_payments"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/get_v1_me"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/get_v1_orgns"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/get_v1_orgns_orgn_id_api_keys"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/get_v1_orgns_orgn_id_members"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/get_v1_orgns_orgn_id_tariffs"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/get_v1_webhook"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/patch_v1_orgns_orgn_id_members_member_id"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/post_v1_auth_email_resend_verification"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/post_v1_auth_email_verify"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/post_v1_auth_login"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/post_v1_auth_password_change"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/post_v1_auth_password_forgot"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/post_v1_auth_password_reset"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/post_v1_auth_refresh"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/post_v1_auth_register"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/post_v1_auth_switch_orgn"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/post_v1_billing_payments"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/post_v1_customer_customer_id_usage"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/post_v1_customers"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/post_v1_customers_customer_id_subscribe"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/post_v1_orgns"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/post_v1_orgns_orgn_id_api_keys"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/post_v1_orgns_orgn_id_tariffs"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/put_v1_orgns_orgn_id"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/openapi/handlers/put_v1_orgns_orgn_id_tariffs_tariff_id"
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
	post_v1_auth_password_forgot.AuthService
	post_v1_auth_password_reset.AuthService
	post_v1_auth_password_change.AuthService
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

	postV1AuthPasswordForgotHandler interface {
		PostV1AuthPasswordForgot(context.Context, openapi.PostV1AuthPasswordForgotRequestObject) (openapi.PostV1AuthPasswordForgotResponseObject, error)
	}

	postV1AuthPasswordResetHandler interface {
		PostV1AuthPasswordReset(context.Context, openapi.PostV1AuthPasswordResetRequestObject) (openapi.PostV1AuthPasswordResetResponseObject, error)
	}

	postV1AuthPasswordChangeHandler interface {
		PostV1AuthPasswordChange(context.Context, openapi.PostV1AuthPasswordChangeRequestObject) (openapi.PostV1AuthPasswordChangeResponseObject, error)
	}

	postV1AuthSwitchOrgnHandler interface {
		PostV1AuthSwitchOrgn(context.Context, openapi.PostV1AuthSwitchOrgnRequestObject) (openapi.PostV1AuthSwitchOrgnResponseObject, error)
	}

	getV1MeHandler interface {
		GetV1Me(context.Context, openapi.GetV1MeRequestObject) (openapi.GetV1MeResponseObject, error)
	}

	getV1BillingPaymentsHandler interface {
		GetV1BillingPayments(ctx context.Context, request openapi.GetV1BillingPaymentsRequestObject) (_ openapi.GetV1BillingPaymentsResponseObject, _ error)
	}

	postV1BillingPaymentsHandler interface {
		PostV1BillingPayments(ctx context.Context, request openapi.PostV1BillingPaymentsRequestObject) (_ openapi.PostV1BillingPaymentsResponseObject, _ error)
	}

	postV1CustomerCustomerIdUsageHandler interface {
		PostV1CustomerCustomerIdUsage(ctx context.Context, request openapi.PostV1CustomerCustomerIdUsageRequestObject) (_ openapi.PostV1CustomerCustomerIdUsageResponseObject, _ error)
	}

	postV1CustomersHandler interface {
		PostV1Customers(ctx context.Context, request openapi.PostV1CustomersRequestObject) (_ openapi.PostV1CustomersResponseObject, _ error)
	}

	postV1CustomersCustomerIdSubscribeHandler interface {
		PostV1CustomersCustomerIdSubscribe(ctx context.Context, request openapi.PostV1CustomersCustomerIdSubscribeRequestObject) (_ openapi.PostV1CustomersCustomerIdSubscribeResponseObject, _ error)
	}

	getV1OrgnsHandler interface {
		GetV1Orgns(ctx context.Context, request openapi.GetV1OrgnsRequestObject) (_ openapi.GetV1OrgnsResponseObject, _ error)
	}

	postV1OrgnsHandler interface {
		PostV1Orgns(ctx context.Context, request openapi.PostV1OrgnsRequestObject) (_ openapi.PostV1OrgnsResponseObject, _ error)
	}

	deleteV1OrgnsOrgnIdHandler interface {
		DeleteV1OrgnsOrgnId(ctx context.Context, request openapi.DeleteV1OrgnsOrgnIdRequestObject) (_ openapi.DeleteV1OrgnsOrgnIdResponseObject, _ error)
	}

	putV1OrgnsOrgnIdHandler interface {
		PutV1OrgnsOrgnId(ctx context.Context, request openapi.PutV1OrgnsOrgnIdRequestObject) (_ openapi.PutV1OrgnsOrgnIdResponseObject, _ error)
	}

	getV1OrgnsOrgnIdApiKeysHandler interface {
		GetV1OrgnsOrgnIdApiKeys(ctx context.Context, request openapi.GetV1OrgnsOrgnIdApiKeysRequestObject) (_ openapi.GetV1OrgnsOrgnIdApiKeysResponseObject, _ error)
	}

	postV1OrgnsOrgnIdApiKeysHandler interface {
		PostV1OrgnsOrgnIdApiKeys(ctx context.Context, request openapi.PostV1OrgnsOrgnIdApiKeysRequestObject) (_ openapi.PostV1OrgnsOrgnIdApiKeysResponseObject, _ error)
	}

	deleteV1OrgnsOrgnIdApiKeysKeyIdHandler interface {
		DeleteV1OrgnsOrgnIdApiKeysKeyId(ctx context.Context, request openapi.DeleteV1OrgnsOrgnIdApiKeysKeyIdRequestObject) (_ openapi.DeleteV1OrgnsOrgnIdApiKeysKeyIdResponseObject, _ error)
	}

	getV1OrgnsOrgnIdMembersHandler interface {
		GetV1OrgnsOrgnIdMembers(ctx context.Context, request openapi.GetV1OrgnsOrgnIdMembersRequestObject) (_ openapi.GetV1OrgnsOrgnIdMembersResponseObject, _ error)
	}

	deleteV1OrgnsOrgnIdMembersMemberIdHandler interface {
		DeleteV1OrgnsOrgnIdMembersMemberId(ctx context.Context, request openapi.DeleteV1OrgnsOrgnIdMembersMemberIdRequestObject) (_ openapi.DeleteV1OrgnsOrgnIdMembersMemberIdResponseObject, _ error)
	}

	patchV1OrgnsOrgnIdMembersMemberIdHandler interface {
		PatchV1OrgnsOrgnIdMembersMemberId(ctx context.Context, request openapi.PatchV1OrgnsOrgnIdMembersMemberIdRequestObject) (_ openapi.PatchV1OrgnsOrgnIdMembersMemberIdResponseObject, _ error)
	}

	getV1OrgnsOrgnIdTariffsHandler interface {
		GetV1OrgnsOrgnIdTariffs(ctx context.Context, request openapi.GetV1OrgnsOrgnIdTariffsRequestObject) (_ openapi.GetV1OrgnsOrgnIdTariffsResponseObject, _ error)
	}

	postV1OrgnsOrgnIdTariffsHandler interface {
		PostV1OrgnsOrgnIdTariffs(ctx context.Context, request openapi.PostV1OrgnsOrgnIdTariffsRequestObject) (_ openapi.PostV1OrgnsOrgnIdTariffsResponseObject, _ error)
	}

	deleteV1OrgnsOrgnIdTariffsTariffIdHandler interface {
		DeleteV1OrgnsOrgnIdTariffsTariffId(ctx context.Context, request openapi.DeleteV1OrgnsOrgnIdTariffsTariffIdRequestObject) (_ openapi.DeleteV1OrgnsOrgnIdTariffsTariffIdResponseObject, _ error)
	}

	putV1OrgnsOrgnIdTariffsTariffIdHandler interface {
		PutV1OrgnsOrgnIdTariffsTariffId(ctx context.Context, request openapi.PutV1OrgnsOrgnIdTariffsTariffIdRequestObject) (_ openapi.PutV1OrgnsOrgnIdTariffsTariffIdResponseObject, _ error)
	}

	getV1WebhookHandler interface {
		GetV1Webhook(ctx context.Context, request openapi.GetV1WebhookRequestObject) (_ openapi.GetV1WebhookResponseObject, _ error)
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
	postV1AuthPasswordForgotHandler
	postV1AuthPasswordResetHandler
	postV1AuthPasswordChangeHandler
	postV1AuthSwitchOrgnHandler
	getV1MeHandler
	getV1BillingPaymentsHandler
	postV1BillingPaymentsHandler
	postV1CustomerCustomerIdUsageHandler
	postV1CustomersHandler
	postV1CustomersCustomerIdSubscribeHandler
	getV1OrgnsHandler
	postV1OrgnsHandler
	deleteV1OrgnsOrgnIdHandler
	putV1OrgnsOrgnIdHandler
	getV1OrgnsOrgnIdApiKeysHandler
	postV1OrgnsOrgnIdApiKeysHandler
	deleteV1OrgnsOrgnIdApiKeysKeyIdHandler
	getV1OrgnsOrgnIdMembersHandler
	deleteV1OrgnsOrgnIdMembersMemberIdHandler
	patchV1OrgnsOrgnIdMembersMemberIdHandler
	getV1OrgnsOrgnIdTariffsHandler
	postV1OrgnsOrgnIdTariffsHandler
	deleteV1OrgnsOrgnIdTariffsTariffIdHandler
	putV1OrgnsOrgnIdTariffsTariffIdHandler
	getV1WebhookHandler
}

func NewHandler(
	logger log.Logger,
	authService AuthService,
) http.Handler {
	logger = logger.Set("pkg", "internal/transport/http/openapi")

	hand := &serverHandler{
		postV1AuthRegisterHandler:                 post_v1_auth_register.NewHandler(authService),
		postV1AuthEmailVerifyHandler:              post_v1_auth_email_verify.NewHandler(authService),
		postV1AuthEmailResendVerificationHandler:  post_v1_auth_email_resend_verification.NewHandler(authService),
		postV1AuthLoginHandler:                    post_v1_auth_login.NewHandler(authService),
		postV1AuthRefreshHandler:                  post_v1_auth_refresh.NewHandler(authService),
		getV1AuthSessionsHandler:                  get_v1_auth_sessions.NewHandler(authService),
		deleteV1AuthSessionsHandler:               delete_v1_auth_sessions.NewHandler(authService),
		deleteV1AuthSessionsSessionIdHandler:      delete_v1_auth_sessions_session_id.NewHandler(authService),
		postV1AuthPasswordForgotHandler:           post_v1_auth_password_forgot.NewHandler(authService),
		postV1AuthPasswordResetHandler:            post_v1_auth_password_reset.NewHandler(authService),
		postV1AuthPasswordChangeHandler:           post_v1_auth_password_change.NewHandler(authService),
		postV1AuthSwitchOrgnHandler:               post_v1_auth_switch_orgn.NewHandler(authService),
		getV1MeHandler:                            get_v1_me.NewHandler(authService),
		getV1BillingPaymentsHandler:               get_v1_billing_payments.NewHandler(),
		postV1BillingPaymentsHandler:              post_v1_billing_payments.NewHandler(),
		postV1CustomerCustomerIdUsageHandler:      post_v1_customer_customer_id_usage.NewHandler(),
		postV1CustomersHandler:                    post_v1_customers.NewHandler(),
		postV1CustomersCustomerIdSubscribeHandler: post_v1_customers_customer_id_subscribe.NewHandler(),
		getV1OrgnsHandler:                         get_v1_orgns.NewHandler(),
		postV1OrgnsHandler:                        post_v1_orgns.NewHandler(),
		deleteV1OrgnsOrgnIdHandler:                delete_v1_orgns_orgn_id.NewHandler(),
		putV1OrgnsOrgnIdHandler:                   put_v1_orgns_orgn_id.NewHandler(),
		getV1OrgnsOrgnIdApiKeysHandler:            get_v1_orgns_orgn_id_api_keys.NewHandler(),
		postV1OrgnsOrgnIdApiKeysHandler:           post_v1_orgns_orgn_id_api_keys.NewHandler(),
		deleteV1OrgnsOrgnIdApiKeysKeyIdHandler:    delete_v1_orgns_orgn_id_api_keys_key_id.NewHandler(),
		getV1OrgnsOrgnIdMembersHandler:            get_v1_orgns_orgn_id_members.NewHandler(),
		deleteV1OrgnsOrgnIdMembersMemberIdHandler: delete_v1_orgns_orgn_id_members_member_id.NewHandler(),
		patchV1OrgnsOrgnIdMembersMemberIdHandler:  patch_v1_orgns_orgn_id_members_member_id.NewHandler(),
		getV1OrgnsOrgnIdTariffsHandler:            get_v1_orgns_orgn_id_tariffs.NewHandler(),
		postV1OrgnsOrgnIdTariffsHandler:           post_v1_orgns_orgn_id_tariffs.NewHandler(),
		deleteV1OrgnsOrgnIdTariffsTariffIdHandler: delete_v1_orgns_orgn_id_tariffs_tariff_id.NewHandler(),
		putV1OrgnsOrgnIdTariffsTariffIdHandler:    put_v1_orgns_orgn_id_tariffs_tariff_id.NewHandler(),
		getV1WebhookHandler:                       get_v1_webhook.NewHandler(),
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
