package domain

import "errors"

var (
	ErrDataExists = errors.New("data already exists")
	ErrUnknown    = errors.New("unknown error")
)
