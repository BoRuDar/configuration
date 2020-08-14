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

func (defaultProvider) Provide(field reflect.StructField, v reflect.Value, _ ...string) error {
	valStr := getDefaultTag(field)
	return SetField(field, v, valStr)
}
