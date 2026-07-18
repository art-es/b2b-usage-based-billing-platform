package jwt

import (
	"errors"
	"time"
)

var ErrInvalidToken = errors.New("invalid token")

type Claims struct {
	SessionID string
	UserID    string
	OrgnID     *string
	ExpiresAt time.Time
}

func NewClaims(sid, uid string, oid *string) *Claims {
	return &Claims{
		SessionID: sid,
		UserID:    uid,
		OrgnID:     oid,
	}
}
