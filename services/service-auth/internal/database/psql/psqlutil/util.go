package psqlutil

import (
	"errors"

	"github.com/lib/pq"
)

func IsUniqueViolationError(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}
