// Package configuration provides ability to initialize your custom configuration struct from: flags, environment variables, `default` tag, files (json, yaml)
package configuration

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

// New creates a new instance of the configurator
func New(
	cfgPtr interface{}, // must be a pointer to a struct
	providers []Provider,
	loggingEnabled bool,
	failIfCannotSet bool,
) (configurator, error) {
	if len(providers) == 0 {
		return configurator{}, errors.New("providers not found")
	}

	if reflect.TypeOf(cfgPtr).Kind() != reflect.Ptr {
		return configurator{}, errors.New("not a pointer to the struct")
	}

	gLoggingEnabled = loggingEnabled
	gFailIfCannotSet = failIfCannotSet
	gLogger = log.Printf

	return configurator{
		config:    cfgPtr,
		providers: providers,
	}, nil
}

type configurator struct {
	config    interface{}
	providers []Provider
}

// InitValues sets values into struct field using given set of providers
// respecting their order: first defined -> first executed
func (c configurator) InitValues() error {
	return c.fillUp(c.config)
}

// SetLogger changes gLogger
func (configurator) SetLogger(l Logger) {
	gLogger = l
}

func (c configurator) fillUp(i interface{}, parentPath ...string) error {
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
			if err := c.fillUp(vField.Addr().Interface(), currentPath...); err != nil {
				return err
			}
			continue
		}

		if tField.Type.Kind() == reflect.Ptr && tField.Type.Elem().Kind() == reflect.Struct {
			vField.Set(reflect.New(tField.Type.Elem()))
			if err := c.fillUp(vField.Interface(), currentPath...); err != nil {
				return err
			}
			continue
		}

		if err := c.applyProviders(tField, vField, currentPath); err != nil {
			return err
		}
	}
	return nil
}

func (c configurator) applyProviders(field reflect.StructField, v reflect.Value, currentPath []string) error {
	logf("configurator: current path: %v", currentPath)

	for _, provider := range c.providers {
		if provider.Provide(field, v, currentPath...) {
			logf("\n")
			return nil
		}
	}
	errMsg := fmt.Sprintf("configurator: field [%s] with tags [%v] cannot be set!", field.Name, field.Tag)
	logf(errMsg)
	if gFailIfCannotSet {
		fatalf(errMsg)
	}
	return errors.New(errMsg)
}
