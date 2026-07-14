package bcrypt

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/hash"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (*Service) Generate(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	return string(hash), err
}

func (*Service) Compare(s string, h string) error {
	err := bcrypt.CompareHashAndPassword([]byte(h), []byte(s))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return hash.ErrMismatch
	}
	return err
}
