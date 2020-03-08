package configuration

import "reflect"

// Provider defines interface for existing and future custom providers
type Provider interface {
	Provide(field reflect.StructField, v reflect.Value, pathToField ...string) bool
}
