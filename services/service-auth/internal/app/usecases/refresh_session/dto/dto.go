package dto

import "errors"

var (
	ErrInvalidToken = errors.New("invalid token")
)

type Response struct {
	AccessToken  string
	RefreshToken string
}
