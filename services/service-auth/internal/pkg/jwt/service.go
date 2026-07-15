//go:generate mockgen -source=service.go -destination=service_mock_test.go -package=$GOPACKAGE
package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	domain "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
)

type timeService interface {
	GetCurrentTime() time.Time
}

type internalClaims struct {
	*jwt.RegisteredClaims
	UserID         string  `json:"uid,omitempty"`
	OrganizationID *string `json:"oid,omitempty"`
}

type Service struct {
	parser *jwt.Parser
	logger log.Logger
}

func NewService(timeService timeService, logger log.Logger) *Service {
	parser := jwt.NewParser(
		jwt.WithTimeFunc(timeService.GetCurrentTime),
		jwt.WithExpirationRequired(),
	)

	logger = logger.Set("pkg", "internal/pkg/jwt")

	return &Service{
		parser: parser,
		logger: logger,
	}
}

func (s *Service) Generate(secret []byte, claims *domain.Claims) (string, error) {
	obj := jwt.NewWithClaims(jwt.SigningMethodHS256, &internalClaims{
		RegisteredClaims: &jwt.RegisteredClaims{
			ID:        claims.SessionID,
			ExpiresAt: jwt.NewNumericDate(claims.ExpiresAt),
		},
		UserID:         claims.UserID,
		OrganizationID: claims.OrganizationID,
	})

	token, err := obj.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return token, nil
}

func (s *Service) Parse(secret []byte, token string) (*domain.Claims, error) {
	obj, err := s.parser.ParseWithClaims(token, &internalClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			s.logger.Log(log.Warning).
				Set("signing_method", fmt.Sprintf("%v", token.Header["alg"])).
				Set("message", "unexpected signing method").
				Write()

			return nil, errors.New("unexpected signing method")
		}

		return secret, nil
	})
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	claims, ok := obj.Claims.(*internalClaims)
	if !ok || !obj.Valid || claims.ID == "" || claims.UserID == "" {
		return nil, domain.ErrInvalidToken
	}

	return &domain.Claims{
		SessionID:      claims.ID,
		UserID:         claims.UserID,
		OrganizationID: claims.OrganizationID,
		ExpiresAt:      claims.ExpiresAt.Time,
	}, nil
}
