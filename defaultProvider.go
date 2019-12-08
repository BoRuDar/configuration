package configuration

import (
	"log"
	"reflect"
)

func NewDefaultProvider() defaultProvider {
	return defaultProvider{}
}

type defaultProvider struct{}

func (defaultProvider) Provide(field reflect.StructField, v reflect.Value) bool {
	valStr := getDefaultTag(field)
	if len(valStr) == 0 {
		log.Println("defaultProvider: ", valStr)
		return false
	}

	setField(field, v, valStr)
	return true
}
