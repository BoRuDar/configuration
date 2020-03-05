package configuration

import (
	"reflect"
)

func NewDefaultProvider() defaultProvider {
	return defaultProvider{}
}

type defaultProvider struct{}

func (defaultProvider) Provide(field reflect.StructField, v reflect.Value, _ ...string) bool {
	valStr := getDefaultTag(field)
	if len(valStr) == 0 {
		logf("defaultProvider: getDefaultTag returns empty value")
		return false
	}

	SetField(field, v, valStr)
	logf("defaultProvider: set [%s] to field [%s] with tags [%v]", valStr, field.Name, field.Tag)
	return true
}
