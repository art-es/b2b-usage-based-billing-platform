package validate

import (
	"fmt"
)

const (
	ErrorExistence   = 1
	ErrorFormat      = 2
	ErrorLength      = 3
	ErrorIncorrect   = 4
	ErrorUnique      = 5
	ErrorNotVerified = 6
)

type Error struct {
	Field   string
	Type    int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf(
		"validation error: field=%q type=%q message=%q",
		e.Field, e.Type, e.Message,
	)
}
