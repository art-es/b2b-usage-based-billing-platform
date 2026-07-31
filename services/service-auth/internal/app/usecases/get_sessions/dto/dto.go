package dto

import (
	"time"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/domains/jwt"
)

type Request struct {
	Auth   *jwt.Claims
	Cursor *string
}

type Response struct {
	Sessions   []*Session
	NextCursor *string
}

type Session struct {
	ID        string
	CreatedAt time.Time
}
