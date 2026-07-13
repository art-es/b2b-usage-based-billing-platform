package login

import (
	"context"
	"time"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/session"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
)

type hashService interface {
	Compare(s string, hash string) error
	Generate(s string) (string, error)
}

type jwtService interface {
	Generate(claims *jwt.Claims) (string, error)
}

type timeService interface {
	GetCurrentTime() time.Time
}

type uuidService interface {
	Generate() string
}

type sessionRepository interface {
	Save(ctx context.Context, sess *session.Session) error
}

type userRepository interface {
	FindByEmail(ctx context.Context, email string) (*user.User, error)
}
