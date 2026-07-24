package dto

type Request struct {
	Email    string
	Password string
}

type Response struct {
	AccessToken  string
	RefreshToken string
}
