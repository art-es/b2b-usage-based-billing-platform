package openapiutil

import (
	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
)

func ToUUID(s string) types.UUID {
	u, _ := uuid.Parse(s)
	return u
}
