package refresh_session

import (
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

func validateToken(token string) error {
	return validate.Validate("token", token, tokenRules...)
}
