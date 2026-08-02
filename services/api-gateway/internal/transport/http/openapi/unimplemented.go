package openapi

import (
	"context"
	"errors"

	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/generated/openapi"
)

var errUnimplemented = errors.New("unimplemented")

// PostV1AuthPasswordChange Change the password of authorized user
// (POST /v1/auth/password/change)
func (h *serverHandler) PostV1AuthPasswordChange(ctx context.Context, request openapi.PostV1AuthPasswordChangeRequestObject) (_ openapi.PostV1AuthPasswordChangeResponseObject, _ error) {
	return nil, errUnimplemented
}

// PostV1AuthPasswordForgot Send email with reset password link
// (POST /v1/auth/password/forgot)
func (h *serverHandler) PostV1AuthPasswordForgot(ctx context.Context, request openapi.PostV1AuthPasswordForgotRequestObject) (_ openapi.PostV1AuthPasswordForgotResponseObject, _ error) {
	return nil, errUnimplemented
}

// PostV1AuthPasswordReset Reset the password
// (POST /v1/auth/password/reset)
func (h *serverHandler) PostV1AuthPasswordReset(ctx context.Context, request openapi.PostV1AuthPasswordResetRequestObject) (_ openapi.PostV1AuthPasswordResetResponseObject, _ error) {
	return nil, errUnimplemented
}

// PostV1AuthSwitchOrgn Switch an organization in session
// (POST /v1/auth/switch-orgn)
func (h *serverHandler) PostV1AuthSwitchOrgn(ctx context.Context, request openapi.PostV1AuthSwitchOrgnRequestObject) (_ openapi.PostV1AuthSwitchOrgnResponseObject, _ error) {
	return nil, errUnimplemented
}

// GetV1BillingPayments Get payment status
// (GET /v1/billing/payments)
func (h *serverHandler) GetV1BillingPayments(ctx context.Context, request openapi.GetV1BillingPaymentsRequestObject) (_ openapi.GetV1BillingPaymentsResponseObject, _ error) {
	return nil, errUnimplemented
}

// PostV1BillingPayments Generate a new payment link
// (POST /v1/billing/payments)
func (h *serverHandler) PostV1BillingPayments(ctx context.Context, request openapi.PostV1BillingPaymentsRequestObject) (_ openapi.PostV1BillingPaymentsResponseObject, _ error) {
	return nil, errUnimplemented
}

// PostV1CustomerCustomerIdUsage Add new usage
// (POST /v1/customer/{customerId}/usage)
func (h *serverHandler) PostV1CustomerCustomerIdUsage(ctx context.Context, request openapi.PostV1CustomerCustomerIdUsageRequestObject) (_ openapi.PostV1CustomerCustomerIdUsageResponseObject, _ error) {
	return nil, errUnimplemented
}

// PostV1Customers Create a new customer
// (POST /v1/customers)
func (h *serverHandler) PostV1Customers(ctx context.Context, request openapi.PostV1CustomersRequestObject) (_ openapi.PostV1CustomersResponseObject, _ error) {
	return nil, errUnimplemented
}

// PostV1CustomersCustomerIdSubscribe Subscribe customer to tariff
// (POST /v1/customers/{customerId}/subscribe)
func (h *serverHandler) PostV1CustomersCustomerIdSubscribe(ctx context.Context, request openapi.PostV1CustomersCustomerIdSubscribeRequestObject) (_ openapi.PostV1CustomersCustomerIdSubscribeResponseObject, _ error) {
	return nil, errUnimplemented
}

// GetV1Orgns Get organizations
// (GET /v1/orgns)
func (h *serverHandler) GetV1Orgns(ctx context.Context, request openapi.GetV1OrgnsRequestObject) (_ openapi.GetV1OrgnsResponseObject, _ error) {
	return nil, errUnimplemented
}

// PostV1Orgns Create a new organization
// (POST /v1/orgns)
func (h *serverHandler) PostV1Orgns(ctx context.Context, request openapi.PostV1OrgnsRequestObject) (_ openapi.PostV1OrgnsResponseObject, _ error) {
	return nil, errUnimplemented
}

// DeleteV1OrgnsOrgnId Delete the organization
// (DELETE /v1/orgns/{orgnId})
func (h *serverHandler) DeleteV1OrgnsOrgnId(ctx context.Context, request openapi.DeleteV1OrgnsOrgnIdRequestObject) (_ openapi.DeleteV1OrgnsOrgnIdResponseObject, _ error) {
	return nil, errUnimplemented
}

// PutV1OrgnsOrgnId Update the organization
// (PUT /v1/orgns/{orgnId})
func (h *serverHandler) PutV1OrgnsOrgnId(ctx context.Context, request openapi.PutV1OrgnsOrgnIdRequestObject) (_ openapi.PutV1OrgnsOrgnIdResponseObject, _ error) {
	return nil, errUnimplemented
}

// GetV1OrgnsOrgnIdApiKeys Get orgn API keys
// (GET /v1/orgns/{orgnId}/api-keys)
func (h *serverHandler) GetV1OrgnsOrgnIdApiKeys(ctx context.Context, request openapi.GetV1OrgnsOrgnIdApiKeysRequestObject) (_ openapi.GetV1OrgnsOrgnIdApiKeysResponseObject, _ error) {
	return nil, errUnimplemented
}

// PostV1OrgnsOrgnIdApiKeys Create a new organization API key
// (POST /v1/orgns/{orgnId}/api-keys)
func (h *serverHandler) PostV1OrgnsOrgnIdApiKeys(ctx context.Context, request openapi.PostV1OrgnsOrgnIdApiKeysRequestObject) (_ openapi.PostV1OrgnsOrgnIdApiKeysResponseObject, _ error) {
	return nil, errUnimplemented
}

// DeleteV1OrgnsOrgnIdApiKeysKeyId Delete organization API key
// (DELETE /v1/orgns/{orgnId}/api-keys/{keyId})
func (h *serverHandler) DeleteV1OrgnsOrgnIdApiKeysKeyId(ctx context.Context, request openapi.DeleteV1OrgnsOrgnIdApiKeysKeyIdRequestObject) (_ openapi.DeleteV1OrgnsOrgnIdApiKeysKeyIdResponseObject, _ error) {
	return nil, errUnimplemented
}

// GetV1OrgnsOrgnIdMembers Get organization members
// (GET /v1/orgns/{orgnId}/members)
func (h *serverHandler) GetV1OrgnsOrgnIdMembers(ctx context.Context, request openapi.GetV1OrgnsOrgnIdMembersRequestObject) (_ openapi.GetV1OrgnsOrgnIdMembersResponseObject, _ error) {
	return nil, errUnimplemented
}

// DeleteV1OrgnsOrgnIdMembersMemberId Delete organization member
// (DELETE /v1/orgns/{orgnId}/members/{memberId})
func (h *serverHandler) DeleteV1OrgnsOrgnIdMembersMemberId(ctx context.Context, request openapi.DeleteV1OrgnsOrgnIdMembersMemberIdRequestObject) (_ openapi.DeleteV1OrgnsOrgnIdMembersMemberIdResponseObject, _ error) {
	return nil, errUnimplemented
}

// PatchV1OrgnsOrgnIdMembersMemberId Update organization member
// (PATCH /v1/orgns/{orgnId}/members/{memberId})
func (h *serverHandler) PatchV1OrgnsOrgnIdMembersMemberId(ctx context.Context, request openapi.PatchV1OrgnsOrgnIdMembersMemberIdRequestObject) (_ openapi.PatchV1OrgnsOrgnIdMembersMemberIdResponseObject, _ error) {
	return nil, errUnimplemented
}

// GetV1OrgnsOrgnIdTariffs Get tariffs
// (GET /v1/orgns/{orgnId}/tariffs)
func (h *serverHandler) GetV1OrgnsOrgnIdTariffs(ctx context.Context, request openapi.GetV1OrgnsOrgnIdTariffsRequestObject) (_ openapi.GetV1OrgnsOrgnIdTariffsResponseObject, _ error) {
	return nil, errUnimplemented
}

// PostV1OrgnsOrgnIdTariffs Create a new tariff
// (POST /v1/orgns/{orgnId}/tariffs)
func (h *serverHandler) PostV1OrgnsOrgnIdTariffs(ctx context.Context, request openapi.PostV1OrgnsOrgnIdTariffsRequestObject) (_ openapi.PostV1OrgnsOrgnIdTariffsResponseObject, _ error) {
	return nil, errUnimplemented
}

// DeleteV1OrgnsOrgnIdTariffsTariffId Delete the tariff
// (DELETE /v1/orgns/{orgnId}/tariffs/{tariffId})
func (h *serverHandler) DeleteV1OrgnsOrgnIdTariffsTariffId(ctx context.Context, request openapi.DeleteV1OrgnsOrgnIdTariffsTariffIdRequestObject) (_ openapi.DeleteV1OrgnsOrgnIdTariffsTariffIdResponseObject, _ error) {
	return nil, errUnimplemented
}

// PutV1OrgnsOrgnIdTariffsTariffId Update the tariff
// (PUT /v1/orgns/{orgnId}/tariffs/{tariffId})
func (h *serverHandler) PutV1OrgnsOrgnIdTariffsTariffId(ctx context.Context, request openapi.PutV1OrgnsOrgnIdTariffsTariffIdRequestObject) (_ openapi.PutV1OrgnsOrgnIdTariffsTariffIdResponseObject, _ error) {
	return nil, errUnimplemented
}

// GetV1Webhook Connect to webhook to receive events
// (GET /v1/webhook)
func (h *serverHandler) GetV1Webhook(ctx context.Context, request openapi.GetV1WebhookRequestObject) (_ openapi.GetV1WebhookResponseObject, _ error) {
	return nil, errUnimplemented
}
