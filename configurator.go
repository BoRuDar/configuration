// Package configuration provides ability to initialize your custom configuration struct from: flags, environment variables, `default` tag, files (json, yaml)
package configuration

import (
	"fmt"
	"log"
	"reflect"
)

// New creates a new instance of the Configurator.
func New[T any](
	providers ...Provider, // providers will be executed in order of their declaration
) *Configurator[T] {
	return &Configurator[T]{
		configPtr:      new(T),
		providers:      providers,
		registeredTags: map[string]struct{}{},
		loggerFn:       log.Printf,
		onErrorFn: func(field reflect.StructField, err error) {
			if err != nil {
				log.Printf("configurator: field [%s] with tags [%v] cannot be set. Last Provider error: %s", field.Name, field.Tag, err)
			}
		},
		loggingEnabled: false,
	}
}

type Configurator[T any] struct {
	configPtr      *T
	providers      []Provider
	registeredTags map[string]struct{}
	onErrorFn      func(field reflect.StructField, err error)
	loggerFn       func(format string, v ...any)
	loggingEnabled bool
}

func (c *Configurator[T]) SetOptions(options ...ConfiguratorOption[T]) *Configurator[T] {
	for _, o := range options {
		o(c)
	}
	return c
}

// InitValues sets values into struct field using given set of providers
// respecting their order: first defined -> first executed
func (c *Configurator[T]) InitValues() (*T, error) {
	if reflect.TypeOf(c.configPtr).Kind() != reflect.Ptr {
		return nil, ErrNotAPointer
	}

	if len(c.providers) == 0 {
		return nil, ErrNoProviders
	}

	for _, p := range c.providers {
		if _, ok := c.registeredTags[p.Name()]; ok {
			return nil, ErrProviderNameCollision
		}
		c.registeredTags[p.Name()] = struct{}{}

		if err := p.Init(c.configPtr); err != nil {
			return nil, fmt.Errorf("cannot init [%s] provider: %w", p.Name(), err)
		}
	}

	c.fillUp(c.configPtr)
	return c.configPtr, nil
}

func (c *Configurator[T]) fillUp(i any) {
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
			c.fillUp(vField.Addr().Interface())
			continue
		}

		if tField.Type.Kind() == reflect.Ptr && tField.Type.Elem().Kind() == reflect.Struct {
			vField.Set(reflect.New(tField.Type.Elem()))
			c.fillUp(vField.Interface())
			continue
		}

		c.applyProviders(tField, vField)
	}
}

func (c *Configurator[T]) applyProviders(field reflect.StructField, v reflect.Value) {
	if !field.IsExported() {
		return
	}

	var lastErr error
	for _, provider := range c.providers {
		if _, found := fetchTagKey(field.Tag)[provider.Tag()]; !found {
			// skip provider if it's not specified in tags
			continue
		}

		if lastErr = provider.Provide(field, v); lastErr == nil {
			return
		}
	}

	c.onErrorFn(field, lastErr)
}

// FromEnvAndDefault is a shortcut for `New(cfg, NewEnvProvider(), NewDefaultProvider()).InitValues()`.
func FromEnvAndDefault[T any]() (*T, error) {
	return New[T](NewEnvProvider(), NewDefaultProvider()).InitValues()
}
