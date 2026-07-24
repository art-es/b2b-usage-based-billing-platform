package validate

import (
	"fmt"
)

const (
	ErrorExistence ErrorType = iota + 1
	ErrorFormat
	ErrorLength
	ErrorIncorrect
	ErrorUnique
	ErrorNotVerified
)

type ErrorType int8

type Error struct {
	Field   string
	Type    ErrorType
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf(
		"validation error: field=%q type=%q message=%q",
		e.Field, e.Type, e.Message,
	)
}
