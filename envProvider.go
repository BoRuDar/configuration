package configuration

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// NewEnvProvider creates provider which sets values from ENV variables (gets variable name from `env` tag)
func NewEnvProvider() envProvider {
	return envProvider{}
}

type envProvider struct{}

func (envProvider) Init(_ interface{}) error {
	return nil
}

func (envProvider) Provide(field reflect.StructField, v reflect.Value, _ ...string) error {
	key := getEnvTag(field)
	if len(key) == 0 {
		// field doesn't have a proper tag
		return fmt.Errorf("envProvider: key is empty")
	}

	valStr, ok := os.LookupEnv(strings.ToUpper(key))
	if !ok || len(valStr) == 0 {
		return fmt.Errorf("envProvider: %w", ErrEmptyValue)
	}

	return SetField(field, v, valStr)
}
