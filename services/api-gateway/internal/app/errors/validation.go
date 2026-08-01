package errors

import (
	"fmt"
	"strings"
)

type ValidationErrors []*ValidationError

type ValidationError struct {
	Field   string
	Code    int
	Message string
}

func (ee ValidationErrors) Error() string {
	if len(ee) == 0 {
		return "empty validation error list"
	}

	mm := make([]string, 0, len(ee))
	for _, e := range ee {
		mm = append(mm, e.Error())
	}

	return fmt.Sprintf("validation error list: %s", strings.Join(mm, "; "))
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf(
		"validation error: field=%q code=%d message=%q",
		e.Field, e.Code, e.Message,
	)
}
