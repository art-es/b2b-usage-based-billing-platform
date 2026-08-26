package validate

import (
	"errors"
	"fmt"
	"strings"
)

type Rule interface {
	Validate(val any) (bool, error)
	FormError() *Error
}

func Validate(field string, val any, rules ...Rule) error {
	for _, rule := range rules {
		ok, err := rule.Validate(val)
		if err != nil {
			return fmt.Errorf("field %s: %w", field, err)
		}

		if !ok {
			vErr := rule.FormError()
			vErr.Field = field
			return vErr
		}
	}

	return nil
}

type Errors []*Error

func (errs Errors) Error() string {
	if len(errs) == 0 {
		return "empty validation error list"
	}

	msgs := make([]string, 0, len(errs))
	for _, err := range errs {
		msgs = append(msgs, err.Error())
	}

	return fmt.Sprintf("validation error list: %s", strings.Join(msgs, "; "))
}

func Join(vErrs Errors, err error) (Errors, bool) {
	if err == nil {
		return vErrs, true
	}

	var vErr *Error
	if !errors.As(err, &vErr) {
		return vErrs, false
	}

	return append(vErrs, vErr), true
}
