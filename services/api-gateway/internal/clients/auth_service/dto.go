package auth_service

import "time"

type RegisterRequest struct {
	Name     string
	Email    string
	Password string
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
}

type RefreshSessionResponse struct {
	AccessToken  string
	RefreshToken string
}

type GetSessionsRequest struct {
	Cursor *string
}

type GetSessionsResponse struct {
	Sessions   []*Session
	NextCursor *string
}

type Session struct {
	ID        string
	CreatedAt time.Time
}

type GetMeResponse struct {
	SessionID string
	Name      string
	Email     string
	Orgn      *GetMeResponseOrgn
}

type GetMeResponseOrgn struct {
	ID   string
	Name string
}
