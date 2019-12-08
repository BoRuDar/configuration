package configuration

import (
	"errors"
	"reflect"
)

type Provider interface {
	Provide(field reflect.StructField, v reflect.Value) bool
}

func FillUp(i interface{}, providers []Provider) error {
	var (
		t = reflect.TypeOf(i)
		v = reflect.ValueOf(i)
	)

	switch t.Kind() {
	case reflect.Ptr:
		t = t.Elem()
		v = v.Elem()
	default:
		return errors.New("not a pointer to the struct")
	}

	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Type.Kind() == reflect.Struct {
			if err := FillUp(v.Field(i).Addr().Interface(), providers); err != nil {
				return err
			}
			continue
		}

		applyProviders(t.Field(i), v.Field(i), providers)
	}
	return nil
}

func applyProviders(field reflect.StructField, v reflect.Value, providers []Provider) {
	for _, provider := range providers {
		if provider.Provide(field, v) {
			return
		}
	}
}

func setField(field reflect.StructField, v reflect.Value, valStr string) {
	if v.Kind() == reflect.Ptr {
		setPtrValue(field.Type.Elem(), v, valStr)
		return
	}
	setValue(field.Type, v, valStr)
}
