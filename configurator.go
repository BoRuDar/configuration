// Package configuration provides ability to initialize your custom configuration struct from: flags, environment variables, `default` tag, files (json, yaml)
package configuration

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

// New creates a new instance of the configurator.
func New(
	cfgPtr interface{}, // must be a pointer to a struct
	providers ...Provider, // providers will be executed in order of their declaration
) (configurator, error) {
	if len(providers) == 0 {
		return configurator{}, errors.New("providers not found")
	}

	if reflect.TypeOf(cfgPtr).Kind() != reflect.Ptr {
		return configurator{}, errors.New("not a pointer to the struct")
	}

	return configurator{
		config:    cfgPtr,
		providers: providers,
		loggerFn:  log.Printf,
		onErrorFn: func(err error) {
			panic(err)
		},
		loggingEnabled: true,
	}, nil
}

type configurator struct {
	config         interface{}
	providers      []Provider
	onErrorFn      func(err error)
	loggerFn       func(format string, v ...interface{})
	loggingEnabled bool
}

// InitValues sets values into struct field using given set of providers
// respecting their order: first defined -> first executed
func (c configurator) InitValues() {
	c.fillUp(c.config)
}

// SetLogger changes logger
func (c *configurator) SetLogger(l func(format string, v ...interface{})) {
	c.loggerFn = l
	return
}

// SetOnErrorFn changes function which is called on errors
func (c *configurator) SetOnErrorFn(fn func(error)) {
	c.onErrorFn = fn
}

func (c configurator) fillUp(i interface{}, parentPath ...string) {
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
			tField      = t.Field(i)
			vField      = v.Field(i)
			currentPath = append(parentPath, tField.Name)
		)

		if tField.Type.Kind() == reflect.Struct {
			c.fillUp(vField.Addr().Interface(), currentPath...)
			continue
		}

		if tField.Type.Kind() == reflect.Ptr && tField.Type.Elem().Kind() == reflect.Struct {
			vField.Set(reflect.New(tField.Type.Elem()))
			c.fillUp(vField.Interface(), currentPath...)
			continue
		}

		c.applyProviders(tField, vField, currentPath)
	}
}

func (c configurator) applyProviders(field reflect.StructField, v reflect.Value, currentPath []string) {
	c.logf("current path: %v", currentPath)

	var err error
	for _, provider := range c.providers {
		if err = provider.Provide(field, v, currentPath...); err != nil {
			c.logf("provider error: %v \n", err)
			continue
		}
		c.logf("\n")
		return
	}

	c.onErrorFn(fmt.Errorf("field [%s] with tags [%v] cannot be set: %v", field.Name, field.Tag, err))
}

func (c configurator) logf(format string, v ...interface{}) {
	if c.loggingEnabled {
		c.loggerFn(format, v...)
	}
}
