package validatefuncs

import (
	"fmt"
	"net/mail"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/validate"
)

type emailRule struct {
	msg string
}

func Email(msg *string) *emailRule {
	r := &emailRule{
		msg: "must be a valid email address",
	}

	if msg != nil {
		r.msg = *msg
	}

	return r
}

func (*emailRule) Validate(val any) (bool, error) {
	s, ok := val.(string)
	if !ok {
		return false, fmt.Errorf("unknown value type: %T", val)
	}

	_, err := mail.ParseAddress(s)
	return err == nil, nil
}

func (*emailRule) FormError() *validate.Error {
	return &validate.Error{
		Code:    validate.ErrorFormat,
		Message: "must be a valid email address",
	}
}
