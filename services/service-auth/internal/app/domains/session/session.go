package session

import "github.com/google/uuid"

var (
	NewSessionID    = uuid.NewString
	NewRefreshToken = uuid.NewString
)

type Session struct {
	// Info from access token
	ID             string
	UserID         string
	OrganizationID *string

	// Keeps in DB
	RefreshTokenHash string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

func NewSession(userID string) *Session {
	return &Session{
		ID:     NewSessionID(),
		UserID: userID,
	}
}
