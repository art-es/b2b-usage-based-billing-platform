package uuid

import "github.com/google/uuid"

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Generate() string {
	return uuid.NewString()
}
