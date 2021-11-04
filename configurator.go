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
) configurator {
	return configurator{
		configPtr: cfgPtr,
		providers: providers,
		loggerFn:  log.Printf,
		onFailToSetField: func(err error) {
			log.Fatal(err)
		},
		loggingEnabled: false,
	}
}

type configurator struct {
	configPtr        interface{}
	providers        []Provider
	onFailToSetField func(err error)
	loggerFn         func(format string, v ...interface{})
	loggingEnabled   bool
}

// InitValues sets values into struct field using given set of providers
// respecting their order: first defined -> first executed
func (c configurator) InitValues() error {
	if reflect.TypeOf(c.configPtr).Kind() != reflect.Ptr {
		return ErrNotAPointer
	}

	if len(c.providers) == 0 {
		return errors.New("providers not found")
	}

	c.fillUp(c.configPtr)
	return nil
}

// SetLogger changes logger
func (c *configurator) SetLogger(l func(format string, v ...interface{})) {
	c.loggerFn = l
	return
}

// EnableLogging changes logger
func (c *configurator) EnableLogging(enable bool) {
	c.loggingEnabled = enable
	return
}

// SetOnFailFn sets function which will be called when no value set into the field
func (c *configurator) SetOnFailFn(fn func(error)) {
	c.onFailToSetField = fn
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
	c.logf("configurator: current path: %v", currentPath)

	for _, provider := range c.providers {
		err := provider.Provide(field, v, currentPath...)
		if err == nil {
			return
		}
		c.logf("configurator: %v", err)
	}

	c.onFailToSetField(fmt.Errorf("configurator: field [%s] with tags [%v] cannot be set", field.Name, field.Tag))
}

func (c configurator) logf(format string, v ...interface{}) {
	if c.loggingEnabled {
		c.loggerFn(format, v...)
	}
}
