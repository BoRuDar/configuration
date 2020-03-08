package configuration

import (
	"os"
	"reflect"
	"strings"
)

// NewEnvProvider creates provider which sets values from ENV variables (gets variable name from `env` tag)
func NewEnvProvider() envProvider {
	return envProvider{}
}

type envProvider struct{}

func (envProvider) Provide(field reflect.StructField, v reflect.Value, _ ...string) bool {
	key := getEnvTag(field)
	if len(key) == 0 {
		// field doesn't have a proper tag
		logf("envProvider: key is empty")
		return false
	}

	valStr, ok := os.LookupEnv(strings.ToUpper(key))
	if !ok || len(valStr) == 0 {
		logf("envProvider: os.LookupEnv returns empty value")
		return false
	}

	SetField(field, v, valStr)
	logf("envProvider: set [%s] to field [%s] with tags [%v]", valStr, field.Name, field.Tag)
	return true
}
