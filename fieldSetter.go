package configuration

import (
	"reflect"
	"strconv"
)

func SetField(field reflect.StructField, v reflect.Value, valStr string) {
	if v.Kind() == reflect.Ptr {
		setPtrValue(field.Type.Elem(), v, valStr)
		return
	}
	setValue(field.Type, v, valStr)
}

func setValue(t reflect.Type, v reflect.Value, val string) {
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
		panic("unsupported type: " + v.Kind().String())
	}
}

func setPtrValue(t reflect.Type, v reflect.Value, val string) {
	switch t.Name() {
	case reflect.Int.String(): // doesn't care about 32bit systems
		if i64, err := strconv.ParseInt(val, 10, 64); err == nil {
			i := int(i64)
			v.Set(reflect.ValueOf(&i))
		}
	case reflect.Int8.String():
		if i64, err := strconv.ParseInt(val, 10, 8); err == nil {
			i8 := int8(i64)
			v.Set(reflect.ValueOf(&i8))
		}
	case reflect.Int16.String():
		if i64, err := strconv.ParseInt(val, 10, 16); err == nil {
			i16 := int16(i64)
			v.Set(reflect.ValueOf(&i16))
		}
	case reflect.Int32.String():
		if i64, err := strconv.ParseInt(val, 10, 32); err == nil {
			i32 := int32(i64)
			v.Set(reflect.ValueOf(&i32))
		}
	case reflect.Int64.String():
		if i64, err := strconv.ParseInt(val, 10, 64); err == nil {
			v.Set(reflect.ValueOf(&i64))
		}

	case reflect.Uint.String(): // doesn't care about 32bit systems
		if ui64, err := strconv.ParseUint(val, 10, 64); err == nil {
			ui := uint(ui64)
			v.Set(reflect.ValueOf(&ui))
		}
	case reflect.Uint8.String():
		if ui64, err := strconv.ParseUint(val, 10, 8); err == nil {
			ui8 := uint8(ui64)
			v.Set(reflect.ValueOf(&ui8))
		}
	case reflect.Uint16.String():
		if ui64, err := strconv.ParseUint(val, 10, 16); err == nil {
			ui16 := uint16(ui64)
			v.Set(reflect.ValueOf(&ui16))
		}
	case reflect.Uint32.String():
		if ui64, err := strconv.ParseUint(val, 10, 32); err == nil {
			ui32 := uint32(ui64)
			v.Set(reflect.ValueOf(&ui32))
		}
	case reflect.Uint64.String():
		if ui64, err := strconv.ParseUint(val, 10, 64); err == nil {
			v.Set(reflect.ValueOf(&ui64))
		}

	case reflect.Float32.String():
		if f64, err := strconv.ParseFloat(val, 32); err == nil {
			f32 := float32(f64)
			v.Set(reflect.ValueOf(&f32))
		}
	case reflect.Float64.String():
		if f64, err := strconv.ParseFloat(val, 64); err == nil {
			v.Set(reflect.ValueOf(&f64))
		}

	case reflect.String.String():
		if len(val) > 0 {
			v.Set(reflect.ValueOf(&val))
		}

	case reflect.Bool.String():
		if b, err := strconv.ParseBool(val); err == nil {
			v.Set(reflect.ValueOf(&b))
		}
	default:
		panic("unsupported type: " + t.Kind().String())
	}
}
