package psqlutil

import (
	"errors"

	"github.com/lib/pq"
)

const (
	UniqueViolationErrorCode = "23505"
)

func IsUniqueViolationError(err error) bool {
	var pqErr *pq.Error

	return errors.As(err, &pqErr) &&
		pqErr.Code == UniqueViolationErrorCode
}
