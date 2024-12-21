package configuration

import (
	"errors"
)

var (
	ErrNoTag                 = errors.New("no tag")
	ErrTagNotUnique          = errors.New("tag is not unique")
	ErrEmptyValue            = errors.New("empty value")
	ErrNotAStruct            = errors.New("not a struct")
	ErrInvalidInput          = errors.New("invalid input")
	ErrNoProviders           = errors.New("no providers")
	ErrProviderNameCollision = errors.New("provider name collision")
	ErrProviderTagCollision  = errors.New("provider tag collision")
)
