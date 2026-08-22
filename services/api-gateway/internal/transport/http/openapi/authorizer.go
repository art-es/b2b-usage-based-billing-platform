package openapi

import (
	"context"

	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/app/auth"
	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/clients/auth_service"
)

type accessTokenParser interface {
	ParseToken(ctx context.Context) (*auth_service.ParseTokenResponse, error)
}

type authorizer struct {
	tokenParser accessTokenParser
}

func newAuthorizer(tokenParser accessTokenParser) *authorizer {
	return &authorizer{tokenParser: tokenParser}
}

func (a *authorizer) Authorize(ctx context.Context) (*auth.Auth, error) {
	res, err := a.tokenParser.ParseToken(ctx)
	if err != nil {
		return nil, err
	}

	return &auth.Auth{
		SessionID: res.SessionID,
		UserID:    res.UserID,
		OrgnID:    res.OrgnID,
	}, nil
}
