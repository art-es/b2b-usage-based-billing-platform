package dto

import "errors"

var (
	ErrEmailInUse = errors.New("email is already in use")
)

type Request struct {
	Name     string
	Email    string
	Password string
}
