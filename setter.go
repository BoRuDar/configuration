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

		if err := getValue(t.Field(i), v.Field(i), providers); err != nil {
			return err
		}
	}
	return nil
}

func getValue(field reflect.StructField, v reflect.Value, providers []valueProvider) error {
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

func setValue(t reflect.Type, v reflect.Value, val string) error {
	switch t.Kind() {
	case reflect.String:
		v.SetString(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, _ := strconv.ParseInt(val, 10, 64)
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, _ := strconv.ParseUint(val, 10, 64)
		v.SetUint(i)
	case reflect.Float32, reflect.Float64:
		f, _ := strconv.ParseFloat(val, 64)
		v.SetFloat(f)
	case reflect.Bool:
		b, _ := strconv.ParseBool(val)
		v.SetBool(b)
	default:
		return errors.New("unsupported type: " + v.Kind().String())
	}
	return nil
}

func setPtrValue(t reflect.Type, v reflect.Value, val string) error {
	switch t.Name() {
	case reflect.Int.String():
		//reflect.Int8.String(),
		//reflect.Int16.String(),
		//reflect.Int32.String(),
		//reflect.Int64.String():
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			ii := int(i)
			v.Set(reflect.ValueOf(&ii))
		}
	case reflect.String.String():
		if len(val) > 0 {
			v.Set(reflect.ValueOf(&val))
		}
	//case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	//	i, _ := strconv.ParseUint(val, 10, 64)
	//	v.SetUint(i)
	//case reflect.Float32, reflect.Float64:
	//	f, _ := strconv.ParseFloat(val, 64)
	//	v.SetFloat(f)
	case reflect.Bool.String():
		if b, err := strconv.ParseBool(val); err == nil {
			v.Set(reflect.ValueOf(&b))
		}
	default:
		return errors.New("unsupported type: " + v.Kind().String())
	}
	return nil
}
