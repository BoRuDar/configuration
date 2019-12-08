package configuration

import (
	"os"
	"reflect"
	"strings"
)

func NewEnvProvider() envProvider {
	return envProvider{}
}

type envProvider struct{}

func (envProvider) Provide(field reflect.StructField, v reflect.Value) bool {
	key := getEnvTag(field)
	if len(key) == 0 { // if "env" is not set try to use regular json tag
		key = strings.ToUpper(getJSONTag(field))
	}
	if len(key) == 0 {
		// field doesn't have a proper tag
		return false
	}

	valStr, ok := os.LookupEnv(strings.ToUpper(key))
	if !ok || len(valStr) == 0 {
		return false
	}

	setField(field, v, valStr)
	return true
}
