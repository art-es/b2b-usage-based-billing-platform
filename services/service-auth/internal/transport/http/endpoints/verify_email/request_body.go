package verify_email

import (
	"github.com/google/uuid"
)

type requestBody struct {
	Token string `json:"name"`
}

func (b *requestBody) validate() (string, int, bool) {
	switch {
	case len(b.Token) == 0:
		return errMsgRequiredToken, errCodeRequiredToken, false
	case b.validateToken():
		return errMsgInvalidToken, errCodeInvalidToken, false
	default:
		return "", 0, true
	}
}

func (b *requestBody) validateToken() bool {
	_, err := uuid.Parse(b.Token)
	return err == nil
}
