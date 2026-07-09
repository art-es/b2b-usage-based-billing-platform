package bcrypt

import "golang.org/x/crypto/bcrypt"

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (*Service) Generate(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	return string(hash), err
}
