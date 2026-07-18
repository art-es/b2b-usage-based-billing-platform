//go:generate mockgen -source=contract.go -destination=contract_mock_test.go -package=$GOPACKAGE
package login

import (
	"context"
	"time"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/session"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
)

type jwtService interface {
	Generate(secret []byte, claims *jwt.Claims) (string, error)
}

type keyedHashService interface {
	Generate(secret []byte, s string) (string, error)
}

type passwordHashService interface {
	Compare(s string, hash string) error
}

type timeService interface {
	GetCurrentTime() time.Time
}

type uuidService interface {
	Generate() string
}

type sessionRepository interface {
	Save(ctx context.Context, ses *session.Session) error
}

type userRepository interface {
	FindByEmail(ctx context.Context, email string) (*user.User, error)
}
