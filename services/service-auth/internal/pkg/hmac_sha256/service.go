package hmac_sha256

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (*Service) Generate(secret []byte, str string) (string, error) {
	mac := hmac.New(sha256.New, secret)

	_, err := mac.Write([]byte(str))
	if err != nil {
		return "", fmt.Errorf("write to hmac sha256: %w", err)
	}

	return hex.EncodeToString(mac.Sum(nil)), nil
}
