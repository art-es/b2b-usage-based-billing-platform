package update

import (
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/usecases/update/dto"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/validate"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/app/validate/validatefuncs"
	"github.com/art-es/b2b-usage-based-billing-platform/services/service-orgn/internal/pkg/ptr"
)

var errOrgnIDNotFound = &validate.Error{
	Field:   "orgn_id",
	Code:    validate.ErrorIncorrect,
	Message: "not found",
}

var errOrgnIDInvalidFormat = &validate.Error{
	Field:   "orgn_id",
	Code:    validate.ErrorFormat,
	Message: "invalid format",
}

var orgnIDRules = []validate.Rule{
	validatefuncs.Required(),
	validatefuncs.UUID(ptr.To(errOrgnIDInvalidFormat.Message)),
}

var nameRules = []validate.Rule{
	validatefuncs.Required(),
	validatefuncs.Length(ptr.To(3), ptr.To(255)),
}

func validateRequest(r *dto.Request) error {
	var vErrs validate.Errors
	var ok bool

	err := validate.Validate("orgn_id", r.OrgnID, orgnIDRules...)
	vErrs, ok = validate.Join(vErrs, err)
	if !ok {
		return err
	}

	err = validate.Validate("name", r.Name, nameRules...)
	vErrs, ok = validate.Join(vErrs, err)
	if !ok {
		return err
	}

	return vErrs
}
