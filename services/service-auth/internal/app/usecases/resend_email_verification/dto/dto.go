package dto

import (
	"errors"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/user"
)

var (
	ErrInvalidEmail  = errors.New("invalid email")
	ErrEmailVerified = user.ErrEmailVerified
)
