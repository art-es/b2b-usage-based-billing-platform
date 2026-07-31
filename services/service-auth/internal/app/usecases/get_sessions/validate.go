package get_sessions

import "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/validate"

var errInvalidCursor = &validate.Error{
	Field:   "next_cursor",
	Code:    validate.ErrorIncorrect,
	Message: "invalid cursor",
}
