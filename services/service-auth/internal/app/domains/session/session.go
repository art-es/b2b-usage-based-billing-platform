package session

import (
	"time"
)

const (
	AccessTokenExpiry  = 15 * time.Minute
	RefreshTokenExpiry = 14 * 24 * time.Hour
)

type Session struct {
	ID                    string
	UserID                string
	RefreshTokenHash      string
	RefreshTokenExpiresAt time.Time
}

func NewSession(
	userID string,
	refreshTokenHash string,
	now time.Time,
) *Session {
	return &Session{
		UserID:                userID,
		RefreshTokenHash:      refreshTokenHash,
		RefreshTokenExpiresAt: now.Add(RefreshTokenExpiry),
	}
}

func (s *Session) Stored() bool {
	return s.ID != ""
}
