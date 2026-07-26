package auth_service

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
