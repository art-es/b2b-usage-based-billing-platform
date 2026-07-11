package session

type Session struct {
	// Info from access token
	ID             string
	UserID         string
	OrganizationID *string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}
