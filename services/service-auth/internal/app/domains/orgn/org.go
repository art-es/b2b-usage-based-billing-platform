package orgn

import "errors"

var ErrNotFound = errors.New("orgn not found")

type Orgn struct {
	ID     string
	Name   string
	UserID string
}
