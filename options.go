package configuration

import "reflect"

// ConfiguratorOption defines Option function for Configuration
type ConfiguratorOption[T any] func(*Configurator[T])

// OnFailFnOpt sets function which will be called when an error occurs during Configurator.applyProviders()
func OnFailFnOpt[T any](fn func(reflect.StructField, error)) ConfiguratorOption[T] {
	return func(c *Configurator[T]) {
		c.onErrorFn = fn
	}
}
