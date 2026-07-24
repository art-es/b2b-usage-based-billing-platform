package validatefuncs

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/app/validate"
)

type uuidRule struct {
	msg string
}

func UUID(msg *string) *uuidRule {
	r := &uuidRule{
		msg: "must be a valid uuid address",
	}

	if msg != nil {
		r.msg = *msg
	}

	return r
}

func (*uuidRule) Validate(val any) (bool, error) {
	s, ok := val.(string)
	if !ok {
		return false, fmt.Errorf("unknown value type: %T", val)
	}

	_, err := uuid.Parse(s)
	return err == nil, nil
}

func (r *uuidRule) FormError() *validate.Error {
	return &validate.Error{
		Type:    validate.ErrorFormat,
		Message: r.msg,
	}
}
