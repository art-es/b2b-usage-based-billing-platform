package time

import "time"

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetCurrentTime() time.Time {
	return time.Now()
}
