package validatefuncs

import (
	"fmt"
	"strings"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/validate"
)

type lengthRule struct {
	min *int
	max *int
	msg string
}

func Length(min, max *int) *lengthRule {
	msgs := make([]string, 0, 2)
	if min != nil {
		msgs = append(msgs, fmt.Sprintf("min length %d", *min))
	}
	if min != nil {
		msgs = append(msgs, fmt.Sprintf("max length %d", *max))
	}

	return &lengthRule{
		min: min,
		max: max,
		msg: strings.Join(msgs, ", "),
	}
}

func (r *lengthRule) Validate(val any) (bool, error) {
	s, ok := val.(string)
	if !ok {
		return false, fmt.Errorf("unknown value type: %T", val)
	}

	n := len(s)
	res := (r.min == nil || n >= *r.min) || (r.max == nil || n <= *r.max)
	return res, nil
}

func (r *lengthRule) FormError() *validate.Error {
	return &validate.Error{
		Code:    validate.ErrorLength,
		Message: r.msg,
	}
}
