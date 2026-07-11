package user

import "errors"

var (
	ErrEmailVerified = errors.New("email is already verified")
	ErrUserNotFound  = errors.New("user not found")
)

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
}

func NewRegisteredUser(
	name string,
	email string,
	passwordHash string,
) *User {
	return &User{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	}
}
