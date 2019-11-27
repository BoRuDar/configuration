package configuration

import (
	"errors"
	"reflect"
	"strconv"
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

		if err := setValue(t.Field(i), v.Field(i), providers); err != nil {
			return err
		}
	}
	return nil
}

func setValue(field reflect.StructField, v reflect.Value, providers []valueProvider) error {
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

	switch v.Kind() {
	case reflect.String:
		v.SetString(valStr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, _ := strconv.ParseInt(valStr, 10, 64)
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, _ := strconv.ParseUint(valStr, 10, 64)
		v.SetUint(i)
	case reflect.Bool:
		b, _ := strconv.ParseBool(valStr)
		v.SetBool(b)
	default:
		return errors.New("unsupported type: " + v.Kind().String())
	}

	return nil
}
