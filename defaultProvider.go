package configuration

import (
	"fmt"
	"reflect"
)

// NewDefaultProvider creates new provider which sets values from `default` tag
func NewDefaultProvider() defaultProvider {
	return defaultProvider{}
}

type defaultProvider struct{}

func (defaultProvider) Init(_ interface{}) error {
	return nil
}

func (defaultProvider) Provide(field reflect.StructField, v reflect.Value, _ ...string) error {
	valStr := getDefaultTag(field)
	if len(valStr) == 0 {
		return fmt.Errorf("defaultProvider: %w", ErrEmptyValue)
	}

	return SetField(field, v, valStr)
}
