package switch_orgn

import "github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/validate"

var errIncorrectOrgnID = &validate.Error{
	Field:   "orgn_id",
	Code:    validate.ErrorIncorrect,
	Message: "incorrect",
}
