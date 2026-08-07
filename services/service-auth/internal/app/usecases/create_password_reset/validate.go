package create_password_reset

import (
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/validate"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/validate/validatefuncs"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/ptr"
)

var errIncorrectEmail = &validate.Error{
	Field:   "email",
	Code:    validate.ErrorIncorrect,
	Message: "user with email not found",
}

var errInvalidEmail = &validate.Error{
	Field:   "email",
	Code:    validate.ErrorFormat,
	Message: "invalid email",
}

var emailRules = []validate.Rule{
	validatefuncs.Required(),
	validatefuncs.Length(nil, ptr.To(100)),
	validatefuncs.Email(ptr.To(errInvalidEmail.Message)),
}

func validateEmail(email string) error {
	return validate.Validate("email", email, emailRules...)
}
