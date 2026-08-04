package validate

import (
	"fmt"
)

const (
	ErrorExistence   = 1001
	ErrorFormat      = 1002
	ErrorLength      = 1003
	ErrorIncorrect   = 1004
	ErrorUnique      = 1005
	ErrorNotVerified = 1006
)

type Error struct {
	Field   string
	Code    int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf(
		"validation error: field=%q code=%q message=%q",
		e.Field, e.Code, e.Message,
	)
}
