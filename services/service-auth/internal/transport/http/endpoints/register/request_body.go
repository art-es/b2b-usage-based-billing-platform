package register

import (
	"net/mail"
	"unicode"
)

type requestBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (b *requestBody) validate() (string, int, bool) {
	switch {
	case len(b.Name) == 0:
		return errMsgRequiredName, errCodeRequiredName, false
	case len(b.Email) == 0:
		return errMsgRequiredEmail, errCodeRequiredEmail, false
	case len(b.Password) == 0:
		return errMsgRequiredPassword, errCodeRequiredPassword, false
	case !b.validateName():
		return errMsgInvalidName, errCodeInvalidName, false
	case !b.validateEmail():
		return errMsgInvalidEmail, errCodeInvalidEmail, false
	case !b.validatePassword():
		return errMsgInvalidPassword, errCodeInvalidPassword, false
	default:
		return "", 0, true
	}
}

func (b *requestBody) validateName() bool {
	return len(b.Name) > 2 && len(b.Name) < 101 && unicode.IsLetter(rune(b.Name[0]))
}

func (b *requestBody) validateEmail() bool {
	if len(b.Email) > 100 {
		return false
	}

	_, err := mail.ParseAddress(b.Email)
	return err == nil
}

func (b *requestBody) validatePassword() bool {
	if len(b.Password) < 8 || len(b.Password) > 64 {
		return false
	}

	var hasDigit, hasLetter bool

	for _, r := range b.Password {
		switch {
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsLetter(r):
			hasLetter = true
		}

		if hasDigit && hasLetter {
			return true
		}
	}

	return hasDigit && hasLetter
}
