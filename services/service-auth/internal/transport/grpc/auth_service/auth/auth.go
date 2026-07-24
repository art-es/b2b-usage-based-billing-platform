package auth

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc/metadata"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/grpc/grpcutil"
)

type jwtService interface {
	Parse(secret []byte, token string) (*jwt.Claims, error)
}

type Authorizer struct {
	jwtService jwtService
	jwtSecret  []byte
	logger     log.Logger
}

func NewAuthorizer(jwtService jwtService, jwtSecret string, logger log.Logger) *Authorizer {
	logger = logger.Set("pkg", "internal/grpc/auth_service/authorizer")

	return &Authorizer{
		jwtService: jwtService,
		jwtSecret:  []byte(jwtSecret),
		logger:     logger,
	}
}

func (a *Authorizer) Authorize(ctx context.Context) (*jwt.Claims, error) {
	token, ok := getTokenFromContext(ctx)
	if !ok {
		return nil, grpcutil.Unauthenticated()
	}

	claims, err := a.jwtService.Parse(a.jwtSecret, token)
	if err != nil {
		if errors.Is(err, jwt.ErrInvalidToken) {
			return nil, grpcutil.Unauthenticated()
		}

		a.logger.Log(log.Error).
			Set("message", "unexpected jwt service parse error").
			Set("error", err.Error()).
			Write()

		return nil, grpcutil.InternalError()
	}

	return claims, nil
}

func getTokenFromContext(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}

	values := md["authorization"]
	if len(values) == 0 {
		return "", false
	}

	return parseBearerToken(values[0])
}

func parseBearerToken(s string) (string, bool) {
	ss := strings.Split(s, " ")

	if len(ss) != 2 {
		return "", false
	}

	if strings.ToLower(ss[0]) != "bearer" {
		return "", false
	}

	return ss[1], true
}
