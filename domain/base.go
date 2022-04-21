package domain

import "errors"

var (
	ErrDataExists   = errors.New("data already exists")
	ErrDataNotFound = errors.New("data not found")
	ErrUnknown      = errors.New("unknown error")
)
