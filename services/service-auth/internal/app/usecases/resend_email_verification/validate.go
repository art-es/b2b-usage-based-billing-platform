package resend_email_verification

import (
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/validate"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/validate/validatefuncs"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/ptr"
)

var errInvalidEmail = &validate.Error{
	Field:   "email",
	Type:    validate.ErrorFormat,
	Message: "invalid email",
}

var errEmailNotVerified = &validate.Error{
	Field:   "email",
	Type:    validate.ErrorNotVerified,
	Message: "email is not verified",
}

var emailRules = []validate.Rule{
	validatefuncs.Required(),
	validatefuncs.Length(nil, ptr.To(100)),
	validatefuncs.Email(ptr.To(errInvalidEmail.Message)),
}

func validateEmail(email string) error {
	return validate.Validate("email", email, emailRules...)
}
