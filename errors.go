package configuration

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyValue  = errors.New("empty value")
	ErrNotAPointer = fmt.Errorf("not a pointer to a struct")
)
