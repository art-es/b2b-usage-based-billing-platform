package resend_email_verification

import "net/mail"

type requestBody struct {
	Email string `json:"email"`
}

func (b *requestBody) validate() (string, int, bool) {
	switch {
	case len(b.Email) == 0:
		return errMsgRequiredEmail, errCodeRequiredEmail, false
	case !b.validateEmail():
		return errMsgInvalidEmail, errCodeInvalidEmail, false
	default:
		return "", 0, true
	}
}

func (b *requestBody) validateEmail() bool {
	if len(b.Email) > 100 {
		return false
	}

	_, err := mail.ParseAddress(b.Email)
	return err == nil
}
