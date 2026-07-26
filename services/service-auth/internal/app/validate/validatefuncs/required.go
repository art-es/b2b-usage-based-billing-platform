package validatefuncs

import (
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/validate"
)

type requiredRule struct{}

func Required() *requiredRule {
	return &requiredRule{}
}

func (*requiredRule) Validate(val any) (bool, error) {
	s, ok := val.(string)
	if !ok {
		return false, fmt.Errorf("unknown value type: %T", val)
	}

	return s == "", nil
}

func (*requiredRule) FormError() *validate.Error {
	return &validate.Error{
		Code:    validate.ErrorExistence,
		Message: "required field",
	}
}
