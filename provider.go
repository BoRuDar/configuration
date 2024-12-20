package configuration

import "reflect"

// Provider defines interface for existing and future custom providers.
type Provider interface {
	// Name of the provider
	Name() string
	// Tag name of the provider which will trigger an execution of the provider.
	Tag() string
	// Init accepts a pointer to the struct for the initialization.
	Init(ptr any) error
	// Provide operates on reflect.StructField and reflect.Value to set appropriate value.
	Provide(field reflect.StructField, v reflect.Value) error
}
