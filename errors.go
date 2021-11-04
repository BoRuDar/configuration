package configuration

import (
	"errors"
)

var (
	ErrEmptyValue  = errors.New("empty value")
	ErrNotAPointer = errors.New("not a pointer to a struct")
	ErrNoProviders = errors.New("no providers")
)
