package configuration

import "reflect"

// Provider defines interface for existing and future custom providers.
type Provider interface {
	Name() string
	Tag() string
	Init(ptr any) error
	Provide(field reflect.StructField, v reflect.Value) error
}
