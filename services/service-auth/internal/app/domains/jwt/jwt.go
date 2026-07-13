package jwt

import "time"

type Claims struct {
	SessionID      string
	UserID         string
	OrganizationID *string
	ExpiresAt      time.Time
}

func NewClaims(sid, uid string) *Claims {
	return &Claims{
		SessionID: sid,
		UserID:    uid,
	}
}
