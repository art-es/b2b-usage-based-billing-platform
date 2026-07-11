package dto

import (
	"errors"
)

var (
	ErrWrongCredentials = errors.New("wrong credentials")
	ErrEmailNotVerified = errors.New("email is not verified")
)

type Request struct {
	Email    string
	Password string
}
