package configuration

import (
	"fmt"
	"reflect"
)

const (
	DefaultProviderName = `DefaultProvider`
	DefaultProviderTag  = `default`
)

// NewDefaultProvider creates new provider which sets values from `default` tag
// nolint:revive
func NewDefaultProvider() defaultProvider {
	return defaultProvider{}
}

type defaultProvider struct{}

func (defaultProvider) Name() string {
	return DefaultProviderName
}

func (defaultProvider) Tag() string {
	return DefaultProviderTag
}

func (defaultProvider) Init(_ any) error {
	return nil
}

func (dp defaultProvider) Provide(field reflect.StructField, v reflect.Value) error {
	valStr := field.Tag.Get(DefaultProviderTag)
	if len(valStr) == 0 {
		return fmt.Errorf("%s: %w", DefaultProviderName, ErrEmptyValue)
	}

	return SetField(field, v, valStr)
}
