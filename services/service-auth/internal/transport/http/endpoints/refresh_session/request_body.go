package refresh_session

import (
	"github.com/google/uuid"
)

type requestBody struct {
	Token string `json:"refresh_token"`
}

func (b *requestBody) validate() (string, int, bool) {
	switch {
	case len(b.Token) == 0:
		return errMsgRequiredToken, errCodeRequiredToken, false
	case !b.validateToken(b.Token):
		return errMsgInvalidToken, errCodeInvalidToken, false
	default:
		return "", 0, true
	}
}

func (b *requestBody) validateToken(s string) bool {
	err := uuid.Validate(s)
	return err == nil
}
