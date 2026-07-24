package login

import (
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/login/dto"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/validate"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/validate/validatefuncs"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/ptr"
)

var errWrongCredentials = &validate.Error{
	Field:   "email|password",
	Type:    validate.ErrorIncorrect,
	Message: "wrong credentials",
}

var errEmailNotVerified = &validate.Error{
	Field:   "email",
	Type:    validate.ErrorNotVerified,
	Message: "email is not verified",
}

var emailRules = []validate.Rule{
	validatefuncs.Required(),
	validatefuncs.Length(nil, ptr.To(100)),
	validatefuncs.Email(nil),
}

var passwordRules = []validate.Rule{
	validatefuncs.Required(),
	validatefuncs.Length(ptr.To(8), ptr.To(64)),
	validatefuncs.Password(),
}

func validateRequest(r *dto.Request) error {
	var vErrs validate.Errors
	var ok bool

	err := validate.Validate("email", r.Email, emailRules...)
	vErrs, ok = validate.Join(vErrs, err)
	if !ok {
		return err
	}

	err = validate.Validate("password", r.Password, passwordRules...)
	vErrs, ok = validate.Join(vErrs, err)
	if !ok {
		return err
	}

	return vErrs
}
