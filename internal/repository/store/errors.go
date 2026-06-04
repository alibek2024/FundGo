package store

import "errors"

var (
	ErrNotFound      = errors.New("entity not found")
	ErrAlreadyExists = errors.New("entity already exists")
	ErrDataConflict  = errors.New("data integrity conflict")
	ErrInternalDB    = errors.New("internal database error")
)
