package configuration

import (
	"errors"
)

var (
	ErrNoTag        = errors.New("no tag")
	ErrTagNotUnique = errors.New("tag is not unique")
	ErrEmptyValue   = errors.New("empty value")
	ErrNotAPointer  = errors.New("not a pointer to a struct")
	ErrNoProviders  = errors.New("no providers")
)
