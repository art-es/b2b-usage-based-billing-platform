//go:generate mockgen -source=contract.go -destination=contract_mock_test.go -package=$GOPACKAGE
package refresh_session

import (
	"context"
	"time"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/session"
)

type jwtService interface {
	Generate(secret []byte, claims *jwt.Claims) (string, error)
}

type keyedHashService interface {
	Generate(secret []byte, s string) (string, error)
}

type timeService interface {
	GetCurrentTime() time.Time
}

type uuidService interface {
	Generate() string
}

type sessionRepository interface {
	FindByRefreshTokenHash(ctx context.Context, hash string) (*session.Session, error)
	Save(ctx context.Context, ses *session.Session) error
}
