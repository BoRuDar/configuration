package configuration

import "reflect"

// Provider defines interface for existing and future custom providers
type Provider interface {
	// Name() string
	// OnErrorFn()
	Init(ptr interface{}) error
	Provide(field reflect.StructField, v reflect.Value, pathToField ...string) error
}
