package openapi

import (
	"net/http"

	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/transport/http/httputil"
)

type handler struct {
	logger log.Logger
}

func NewHandler(logger log.Logger) http.Handler {
	logger = logger.Set("pkg", "internal/transport/http/handler")

	return openapi.HandlerWithOptions(
		&handler{
			logger: logger,
		},
		openapi.StdHTTPServerOptions{
			ErrorHandlerFunc: httputil.HandleError,
		},
	)
}

// PostV1AuthEmailResendVerification Resend email verification
// (POST /v1/auth/email/resend-verification)
func (h *handler) PostV1AuthEmailResendVerification(w http.ResponseWriter, r *http.Request) {
	httputil.WriteNotImplemented(w, h.logger, "POST /v1/auth/email/resend-verification")
}

// PostV1AuthEmailVerify Verify email
// (POST /v1/auth/email/verify)
func (h *handler) PostV1AuthEmailVerify(w http.ResponseWriter, r *http.Request) {
	httputil.WriteNotImplemented(w, h.logger, "POST /v1/auth/email/verify")
}

// PostV1AuthLogin Create a new session
// (POST /v1/auth/login)
func (h *handler) PostV1AuthLogin(w http.ResponseWriter, r *http.Request) {
	httputil.WriteNotImplemented(w, h.logger, "POST /v1/auth/login")
}

// PostV1AuthPasswordChange Change the password of authorized user
// (POST /v1/auth/password/change)
func (h *handler) PostV1AuthPasswordChange(w http.ResponseWriter, r *http.Request) {
	httputil.WriteNotImplemented(w, h.logger, "POST /v1/auth/password/change")
}

// PostV1AuthPasswordForgot Send email with reset password link
// (POST /v1/auth/password/forgot)
func (h *handler) PostV1AuthPasswordForgot(w http.ResponseWriter, r *http.Request) {
	httputil.WriteNotImplemented(w, h.logger, "POST /v1/auth/password/forgot")
}

// PostV1AuthPasswordReset Reset the password
// (POST /v1/auth/password/reset)
func (h *handler) PostV1AuthPasswordReset(w http.ResponseWriter, r *http.Request) {
	httputil.WriteNotImplemented(w, h.logger, "POST /v1/auth/password/reset")
}

// PostV1AuthRefresh Refresh an access token of session
// (POST /v1/auth/refresh)
func (h *handler) PostV1AuthRefresh(w http.ResponseWriter, r *http.Request) {
	httputil.WriteNotImplemented(w, h.logger, "POST /v1/auth/refresh")
}

// PostV1AuthRegister Create a new account and send email verification
// (POST /v1/auth/register)
func (h *handler) PostV1AuthRegister(w http.ResponseWriter, r *http.Request) {
	httputil.WriteNotImplemented(w, h.logger, "POST /v1/auth/register")
}

// DeleteV1AuthSessions Finish all sessions
// (DELETE /v1/auth/sessions)
func (h *handler) DeleteV1AuthSessions(w http.ResponseWriter, r *http.Request) {
	httputil.WriteNotImplemented(w, h.logger, "DELETE /v1/auth/sessions")
}

// GetV1AuthSessions Get sessions
// (GET /v1/auth/sessions)
func (h *handler) GetV1AuthSessions(w http.ResponseWriter, r *http.Request) {
	httputil.WriteNotImplemented(w, h.logger, "GET /v1/auth/sessions")
}

// DeleteV1AuthSessionsSessionId Finish the session
// (DELETE /v1/auth/sessions/{sessionId})
func (h *handler) DeleteV1AuthSessionsSessionId(w http.ResponseWriter, r *http.Request, sessionId string) {
	httputil.WriteNotImplemented(w, h.logger, "DELETE /v1/auth/sessions/{sessionId}")
}

// PostV1AuthSwitchOrgn Switch an organization in session
// (POST /v1/auth/switch-orgn)
func (h *handler) PostV1AuthSwitchOrgn(w http.ResponseWriter, r *http.Request) {
	httputil.WriteNotImplemented(w, h.logger, "POST /v1/auth/switch-orgn")
}

// GetV1BillingPayments Get payment status
// (GET /v1/billing/payments)
func (h *handler) GetV1BillingPayments(w http.ResponseWriter, r *http.Request) {
	httputil.WriteNotImplemented(w, h.logger, "GET /v1/billing/payments")
}

// PostV1BillingPayments Generate a new payment link
// (POST /v1/billing/payments)
func (h *handler) PostV1BillingPayments(w http.ResponseWriter, r *http.Request) {
	httputil.WriteNotImplemented(w, h.logger, "POST /v1/billing/payments")
}

// PostV1CustomerCustomerIdUsage Add new usage
// (POST /v1/customer/{customerId}/usage)
func (h *handler) PostV1CustomerCustomerIdUsage(w http.ResponseWriter, r *http.Request, customerId string) {
	httputil.WriteNotImplemented(w, h.logger, "POST /v1/customer/{customerId}/usage")
}

// PostV1Customers Create a new customer
// (POST /v1/customers)
func (h *handler) PostV1Customers(w http.ResponseWriter, r *http.Request) {
	httputil.WriteNotImplemented(w, h.logger, "POST /v1/customers")
}

// PostV1CustomersCustomerIdSubscribe Subscribe customer to tariff
// (POST /v1/customers/{customerId}/subscribe)
func (h *handler) PostV1CustomersCustomerIdSubscribe(w http.ResponseWriter, r *http.Request, customerId string) {
	httputil.WriteNotImplemented(w, h.logger, "POST /v1/customers/{customerId}/subscribe")
}

// GetV1Me Get user info of current session
// (GET /v1/me)
func (h *handler) GetV1Me(w http.ResponseWriter, r *http.Request) {
	httputil.WriteNotImplemented(w, h.logger, "GET /v1/me")
}

// GetV1Orgns Get organizations
// (GET /v1/orgns)
func (h *handler) GetV1Orgns(w http.ResponseWriter, r *http.Request, params openapi.GetV1OrgnsParams) {
	httputil.WriteNotImplemented(w, h.logger, "GET /v1/orgns")
}

// PostV1Orgns Create a new organization
// (POST /v1/orgns)
func (h *handler) PostV1Orgns(w http.ResponseWriter, r *http.Request) {
	httputil.WriteNotImplemented(w, h.logger, "POST /v1/orgns")
}

// DeleteV1OrgnsOrgnId Delete the organization
// (DELETE /v1/orgns/{orgnId})
func (h *handler) DeleteV1OrgnsOrgnId(w http.ResponseWriter, r *http.Request, orgnId string) {
	httputil.WriteNotImplemented(w, h.logger, "DELETE /v1/orgns/{orgnId}")
}

// PutV1OrgnsOrgnId Update the organization
// (PUT /v1/orgns/{orgnId})
func (h *handler) PutV1OrgnsOrgnId(w http.ResponseWriter, r *http.Request, orgnId string) {
	httputil.WriteNotImplemented(w, h.logger, "PUT /v1/orgns/{orgnId}")
}

// GetV1OrgnsOrgnIdApiKeys Get orgn API keys
// (GET /v1/orgns/{orgnId}/api-keys)
func (h *handler) GetV1OrgnsOrgnIdApiKeys(w http.ResponseWriter, r *http.Request, orgnId string) {
	httputil.WriteNotImplemented(w, h.logger, "GET /v1/orgns/{orgnId}/api-keys")
}

// PostV1OrgnsOrgnIdApiKeys Create a new organization API key
// (POST /v1/orgns/{orgnId}/api-keys)
func (h *handler) PostV1OrgnsOrgnIdApiKeys(w http.ResponseWriter, r *http.Request, orgnId string) {
	httputil.WriteNotImplemented(w, h.logger, "POST /v1/orgns/{orgnId}/api-keys")
}

// DeleteV1OrgnsOrgnIdApiKeysKeyId Delete organization API key
// (DELETE /v1/orgns/{orgnId}/api-keys/{keyId})
func (h *handler) DeleteV1OrgnsOrgnIdApiKeysKeyId(w http.ResponseWriter, r *http.Request, orgnId string, keyId string) {
	httputil.WriteNotImplemented(w, h.logger, "DELETE /v1/orgns/{orgnId}/api-keys/{keyId}")
}

// GetV1OrgnsOrgnIdMembers Get organization members
// (GET /v1/orgns/{orgnId}/members)
func (h *handler) GetV1OrgnsOrgnIdMembers(w http.ResponseWriter, r *http.Request, orgnId string) {
	httputil.WriteNotImplemented(w, h.logger, "GET /v1/orgns/{orgnId}/members")
}

// DeleteV1OrgnsOrgnIdMembersMemberId Delete organization member
// (DELETE /v1/orgns/{orgnId}/members/{memberId})
func (h *handler) DeleteV1OrgnsOrgnIdMembersMemberId(w http.ResponseWriter, r *http.Request, orgnId string, memberId string) {
	httputil.WriteNotImplemented(w, h.logger, "DELETE /v1/orgns/{orgnId}/members/{memberId}")
}

// PatchV1OrgnsOrgnIdMembersMemberId Update organization member
// (PATCH /v1/orgns/{orgnId}/members/{memberId})
func (h *handler) PatchV1OrgnsOrgnIdMembersMemberId(w http.ResponseWriter, r *http.Request, orgnId string, memberId string) {
	httputil.WriteNotImplemented(w, h.logger, "PATCH /v1/orgns/{orgnId}/members/{memberId}")
}

// GetV1OrgnsOrgnIdTariffs Get tariffs
// (GET /v1/orgns/{orgnId}/tariffs)
func (h *handler) GetV1OrgnsOrgnIdTariffs(w http.ResponseWriter, r *http.Request, orgnId string) {
	httputil.WriteNotImplemented(w, h.logger, "GET /v1/orgns/{orgnId}/tariffs")
}

// PostV1OrgnsOrgnIdTariffs Create a new tariff
// (POST /v1/orgns/{orgnId}/tariffs)
func (h *handler) PostV1OrgnsOrgnIdTariffs(w http.ResponseWriter, r *http.Request, orgnId string) {
	httputil.WriteNotImplemented(w, h.logger, "POST /v1/orgns/{orgnId}/tariffs")
}

// DeleteV1OrgnsOrgnIdTariffsTariffId Delete the tariff
// (DELETE /v1/orgns/{orgnId}/tariffs/{tariffId})
func (h *handler) DeleteV1OrgnsOrgnIdTariffsTariffId(w http.ResponseWriter, r *http.Request, orgnId string, tariffId string) {
	httputil.WriteNotImplemented(w, h.logger, "DELETE /v1/orgns/{orgnId}/tariffs/{tariffId}")
}

// PutV1OrgnsOrgnIdTariffsTariffId Update the tariff
// (PUT /v1/orgns/{orgnId}/tariffs/{tariffId})
func (h *handler) PutV1OrgnsOrgnIdTariffsTariffId(w http.ResponseWriter, r *http.Request, orgnId string, tariffId string) {
	httputil.WriteNotImplemented(w, h.logger, "PUT /v1/orgns/{orgnId}/tariffs/{tariffId}")
}

// GetV1Webhook Connect to webhook to receive events
// (GET /v1/webhook)
func (h *handler) GetV1Webhook(w http.ResponseWriter, r *http.Request) {
	httputil.WriteNotImplemented(w, h.logger, "GET /v1/webhook")
}
