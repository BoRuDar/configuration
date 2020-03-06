package configuration

import "reflect"

type Provider interface {
	Provide(field reflect.StructField, v reflect.Value, pathToField ...string) bool
}
