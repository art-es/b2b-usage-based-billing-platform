package finish_session

import "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/validate"

var errIncorrectSessionId = &validate.Error{
	Field:   "session_id",
	Code:    validate.ErrorIncorrect,
	Message: "incorrect",
}
