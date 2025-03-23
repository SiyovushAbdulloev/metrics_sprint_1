package error

import (
	"errors"
)

var (
	ErrInvalidType        = errors.New("invalid type")
	ErrInvalidValue       = errors.New("invalid value")
	ErrNotFound           = errors.New("not found")
	ErrSomethingWentWrong = errors.New("something went wrong")
)
