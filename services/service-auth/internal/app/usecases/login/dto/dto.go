package dto

import (
	"errors"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/session"
)

var (
	ErrWrongCredentials = errors.New("wrong credentials")
	ErrEmailNotVerified = errors.New("email is not verified")
)

type Request struct {
	Email    string
	Password string
}

type Response session.Tokens
