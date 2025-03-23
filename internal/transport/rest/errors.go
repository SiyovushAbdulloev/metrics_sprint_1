package rest

import "errors"

var (
	errInvalidType  = errors.New("invalid type")
	errInvalidValue = errors.New("invalid value")
	errNotFound     = errors.New("not found")
)
