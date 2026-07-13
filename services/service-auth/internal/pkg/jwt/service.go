package jwt

type Service struct{}

func (s *Service) Encode(secret []byte, claims any) (string, error) {
	return "", nil
}

func (s *Service) Decode(secret []byte, token string, target any) error {
	return nil
}
