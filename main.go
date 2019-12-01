package configuration

import (
	"errors"
	"reflect"
)

func FillUp(i interface{}) error {
	var (
		t         = reflect.TypeOf(i)
		v         = reflect.ValueOf(i)
		providers = []valueProvider{provideFromEnv}
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
			if err := FillUp(v.Field(i).Addr().Interface()); err != nil {
				return err
			}
			continue
		}

		if err := setField(t.Field(i), v.Field(i), providers); err != nil {
			return err
		}
	}
	return nil
}

func setField(field reflect.StructField, v reflect.Value, providers []valueProvider) error {
	var valStr string
	for _, fn := range providers {
		valStr, _ = fn(getJSONTag(field)).(string)
		if len(valStr) > 0 {
			break
		}
	}
	if len(valStr) == 0 {
		valStr = getDefaultTag(field)
	}

	if v.Kind() == reflect.Ptr {
		return setPtrValue(field.Type.Elem(), v, valStr)
	}
	return setValue(field.Type, v, valStr)
}
