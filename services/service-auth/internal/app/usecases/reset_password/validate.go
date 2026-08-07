package reset_password

import (
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/usecases/reset_password/dto"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/validate"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/validate/validatefuncs"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/ptr"
)

var errInvalidToken = &validate.Error{
	Field:   "token",
	Code:    validate.ErrorFormat,
	Message: "invalid token",
}

var tokenRules = []validate.Rule{
	validatefuncs.Required(),
	validatefuncs.UUID(ptr.To(errInvalidToken.Message)),
}

var passwordRules = []validate.Rule{
	validatefuncs.Required(),
	validatefuncs.Length(ptr.To(8), ptr.To(64)),
	validatefuncs.Password(),
}

func validateRequest(r *dto.Request) error {
	var vErrs validate.Errors
	var ok bool

	err := validate.Validate("token", r.Token, tokenRules...)
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
