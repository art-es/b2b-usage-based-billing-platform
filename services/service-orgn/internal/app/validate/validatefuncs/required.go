package validatefuncs

import (
	"fmt"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/validate"
)

type requiredRule struct{}

func Required() *requiredRule {
	return &requiredRule{}
}

func (*requiredRule) Validate(val any) (bool, error) {
	switch val := val.(type) {
	case string:
		return len(val) != 0, nil
	case *string:
		return val != nil && len(*val) != 0, nil
	default:
		return false, fmt.Errorf("unknown value type: %T", val)
	}
}

func (*requiredRule) FormError() *validate.Error {
	return &validate.Error{
		Code:    validate.ErrorExistence,
		Message: "required field",
	}
}
