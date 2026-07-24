package auth_service

const (
	LoginResponseCodeOK = iota
	LoginResponseCodeWrongCredentials
	LoginResponseCodeEmailNotVerified
)

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	RequestError *RequestError

	AccessToken  string
	RefreshToken string
}

type RequestError struct {
	Code    int
	Message string
}
