package configuration

import (
	"errors"
	"reflect"
)

func New(
	cfgPtr interface{},
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

	return configurator{
		config:    cfgPtr,
		providers: providers,
	}, nil
}

type configurator struct {
	config    interface{}
	providers []Provider
}

func (c configurator) InitValues() error {
	return c.fillUp(c.config)
}

func (c configurator) fillUp(i interface{}) error {
	var (
		t = reflect.TypeOf(i)
		v = reflect.ValueOf(i)
	)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Type.Kind() == reflect.Struct {
			if err := c.fillUp(v.Field(i).Addr().Interface()); err != nil {
				return err
			}
			continue
		}

		c.applyProviders(t.Field(i), v.Field(i))
	}
	return nil
}

func (c configurator) applyProviders(field reflect.StructField, v reflect.Value) {
	for _, provider := range c.providers {
		if provider.Provide(field, v) {
			logf("\n")
			return
		}
	}
	logf("configurator: field [%s] with tags [%v] cannot be set!", field.Name, field.Tag)
	fail("configurator: field [%s] with tags [%v] cannot be set!", field.Name, field.Tag)
}
