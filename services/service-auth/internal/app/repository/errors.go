package repository

import "errors"

var (
	ErrUnique   = errors.New("unique error")
	ErrNotFound = errors.New("not found")
)
