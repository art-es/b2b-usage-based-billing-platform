package validatefuncs

import (
	"fmt"
	"unicode"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/validate"
)

type passwordRule struct{}

func Password() *passwordRule {
	return &passwordRule{}
}

func (*passwordRule) Validate(val any) (bool, error) {
	s, ok := val.(string)
	if !ok {
		return false, fmt.Errorf("unknown value type: %T", val)
	}

	var hasDigit, hasLetter bool

	for _, r := range s {
		switch {
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsLetter(r):
			hasLetter = true
		}

		if hasDigit && hasLetter {
			return true, nil
		}
	}

	return hasDigit && hasLetter, nil
}

func (*passwordRule) FormError() *validate.Error {
	return &validate.Error{
		Code:    validate.ErrorFormat,
		Message: "must contain at least 1 letter and 1 digit",
	}
}
