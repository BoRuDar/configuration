package configuration

import "reflect"

type Provider interface {
	Provide(field reflect.StructField, v reflect.Value) bool
}
