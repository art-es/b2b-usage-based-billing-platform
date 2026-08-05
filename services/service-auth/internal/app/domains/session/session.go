package session

import (
	"time"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/ptr"
)

const (
	AccessTokenExpiry  = 15 * time.Minute
	RefreshTokenExpiry = 14 * 24 * time.Hour
)

type Session struct {
	ID                    string
	UserID                string
	OrganizationID        *string
	RefreshTokenHash      string
	RefreshTokenExpiresAt time.Time
	CreatedAt             time.Time
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

func (s *Session) SetRefreshTokenHash(hash string, now time.Time) {
	s.RefreshTokenHash = hash
	s.RefreshTokenExpiresAt = now.Add(RefreshTokenExpiry)
}

func (s *Session) SetOrgn(id string) {
	s.OrganizationID = ptr.To(id)
}
