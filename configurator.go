// Package configuration provides ability to initialize your custom configuration struct from: flags, environment variables, `default` tag, files (json, yaml)
package configuration

import (
	"fmt"
	"reflect"
)

// New creates a new instance of the Configurator.
func New[T any](
	providers ...Provider, // providers will be executed in order of their declaration
) *Configurator[T] {
	return &Configurator[T]{
		configPtr:           new(T),
		providers:           providers,
		registeredProviders: map[string]struct{}{},
	}
}

type Configurator[T any] struct {
	configPtr           *T
	providers           []Provider
	registeredProviders map[string]struct{}
}

// InitValues sets values into struct field using given set of providers
// respecting their order: first defined -> first executed
func (c *Configurator[T]) InitValues() (*T, error) {
	if reflect.TypeOf(c.configPtr).Elem().Kind() != reflect.Struct {
		return nil, ErrNotAStruct
	}

	if len(c.providers) == 0 {
		return nil, ErrNoProviders
	}

	for _, p := range c.providers {
		if _, ok := c.registeredProviders[p.Name()]; ok {
			return nil, ErrProviderNameCollision
		}
		c.registeredProviders[p.Name()] = struct{}{}

		if err := p.Init(c.configPtr); err != nil {
			return nil, fmt.Errorf("cannot init [%s] provider: %w", p.Name(), err)
		}
	}

	if err := c.fillUp(c.configPtr); err != nil {
		return nil, err
	}

	return c.configPtr, nil
}

func (c *Configurator[T]) fillUp(i any) error {
	var (
		t = reflect.TypeOf(i)
		v = reflect.ValueOf(i)
	)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		var (
			tField = t.Field(i)
			vField = v.Field(i)
		)

		if tField.Type.Kind() == reflect.Struct {
			if err := c.fillUp(vField.Addr().Interface()); err != nil {
				return err
			}
			continue
		}

		if tField.Type.Kind() == reflect.Ptr && tField.Type.Elem().Kind() == reflect.Struct {
			vField.Set(reflect.New(tField.Type.Elem()))
			if err := c.fillUp(vField.Interface()); err != nil {
				return err
			}
			continue
		}

		if err := c.applyProviders(tField, vField); err != nil {
			return err
		}
	}

	return nil
}

func (c *Configurator[T]) applyProviders(field reflect.StructField, v reflect.Value) error {
	if !field.IsExported() {
		return nil
	}

	for _, provider := range c.providers {
		if _, found := fetchTagKey(field.Tag)[provider.Tag()]; !found {
			// skip provider if it's not specified in tags
			continue
		}

		if err := provider.Provide(field, v); err == nil {
			return nil
		}
	}

	return fmt.Errorf("filed [%s] with tags [%s] hasn't been set", field.Name, field.Tag)
}

// FromEnvAndDefault is a shortcut for `New(cfg, NewEnvProvider(), NewDefaultProvider()).InitValues()`.
func FromEnvAndDefault[T any]() (*T, error) {
	return New[T](NewEnvProvider(), NewDefaultProvider()).InitValues()
}
