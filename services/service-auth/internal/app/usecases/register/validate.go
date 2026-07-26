package register

import (
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/register/dto"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/validate"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/validate/validatefuncs"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/ptr"
)

var errEmailInUse = &validate.Error{
	Field:   "email",
	Code:    validate.ErrorUnique,
	Message: "email is already in use",
}

var nameRules = []validate.Rule{
	validatefuncs.Required(),
	validatefuncs.Length(ptr.To(2), ptr.To(101)),
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

	err := validate.Validate("name", r.Name, nameRules...)
	vErrs, ok = validate.Join(vErrs, err)
	if !ok {
		return err
	}

	err = validate.Validate("email", r.Email, emailRules...)
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
