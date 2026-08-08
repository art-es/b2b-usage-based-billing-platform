package change_password

import (
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/change_password/dto"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/validate"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/validate/validatefuncs"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/ptr"
)

var errIncorrectOldPassword = &validate.Error{
	Field:   "old_password",
	Code:    validate.ErrorIncorrect,
	Message: "incorrect old password",
}

var passwordRules = []validate.Rule{
	validatefuncs.Required(),
	validatefuncs.Length(ptr.To(8), ptr.To(64)),
	validatefuncs.Password(),
}

func validateRequest(r *dto.Request) error {
	var vErrs validate.Errors
	var ok bool

	err := validate.Validate("old_password", r.OldPassword, passwordRules...)
	vErrs, ok = validate.Join(vErrs, err)
	if !ok {
		return err
	}

	err = validate.Validate("new_password", r.NewPassword, passwordRules...)
	vErrs, ok = validate.Join(vErrs, err)
	if !ok {
		return err
	}

	return vErrs
}
