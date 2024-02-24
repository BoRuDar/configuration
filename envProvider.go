package configuration

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

const EnvProviderName = `EnvProvider`

// NewEnvProvider creates provider which sets values from ENV variables (gets variable name from `env` tag)
// nolint:revive
func NewEnvProvider() envProvider {
	return envProvider{}
}

type envProvider struct{}

func (envProvider) Name() string {
	return EnvProviderName
}

func (envProvider) Init(_ any) error {
	return nil
}

func (ep envProvider) Provide(field reflect.StructField, v reflect.Value) error {
	key := field.Tag.Get("env")
	if len(key) == 0 {
		// field doesn't have a proper tag
		return fmt.Errorf("%s: key is empty", EnvProviderName)
	}

	valStr, ok := os.LookupEnv(strings.ToUpper(key))
	if !ok || len(valStr) == 0 {
		return fmt.Errorf("%s: %w", EnvProviderName, ErrEmptyValue)
	}

	return SetField(field, v, valStr)
}
