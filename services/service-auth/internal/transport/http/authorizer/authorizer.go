package authorizer

import (
	"errors"
	"net/http"
	"strings"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/transport/http/httputil"
)

type jwtService interface {
	Parse(secret []byte, token string) (*jwt.Claims, error)
}

type Authorizer struct {
	jwtService jwtService
	jwtSecret  []byte
	logger     log.Logger
}

func New(jwtService jwtService, jwtSecret string, logger log.Logger) *Authorizer {
	logger = logger.Set("pkg", "internal/http/authorizer")

	return &Authorizer{
		jwtService: jwtService,
		jwtSecret:  []byte(jwtSecret),
		logger:     logger,
	}
}

func (a *Authorizer) Authorize(w http.ResponseWriter, r *http.Request) (*jwt.Claims, bool) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		httputil.WriteUnauthorized(w)
		return nil, false
	}

	token, ok := getBearerToken(authHeader)
	if !ok {
		httputil.WriteUnauthorized(w)
		return nil, false
	}

	claims, err := a.jwtService.Parse(a.jwtSecret, token)
	if err != nil {
		if errors.Is(err, jwt.ErrInvalidToken) {
			httputil.WriteUnauthorized(w)
			return nil, false
		}

		a.logger.Log(log.Error).
			Set("message", "unexpected jwt service parse error").
			Set("error", err.Error()).
			Write()

		httputil.WriteInternalError(w)
		return nil, false
	}

	return claims, true
}

func getBearerToken(s string) (string, bool) {
	ss := strings.Split(s, " ")

	if len(ss) != 2 {
		return "", false
	}

	if strings.ToLower(ss[0]) != "bearer" {
		return "", false
	}

	return ss[1], true
}
