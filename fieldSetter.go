package configuration

import (
	"reflect"
	"strconv"
	"strings"
)

const sliceSeparator = ";"

// SetField sets field with `valStr` value (converts to the proper type beforehand)
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

	case reflect.Slice:
		setSlice(t, v, val)

	default:
		failf("unsupported type: %v", v.Kind().String())
	}
}

func setSlice(t reflect.Type, v reflect.Value, val string) {
	var items []string
	for _, item := range strings.Split(val, sliceSeparator) {
		item = strings.TrimSpace(item)
		if len(item) > 0 {
			items = append(items, item)
		}
	}

	size := len(items)
	if size < 2 {
		return
	}
	slice := reflect.MakeSlice(t, size, size)

	switch t.Elem().Kind() {
	case reflect.String:
		for i := 0; i < size; i++ {
			slice.Index(i).SetString(items[i])
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		for i := 0; i < size; i++ {
			val, _ := strconv.ParseInt(items[i], 10, 64)
			slice.Index(i).SetInt(val)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		for i := 0; i < size; i++ {
			val, _ := strconv.ParseUint(items[i], 10, 64)
			slice.Index(i).SetUint(val)
		}
	case reflect.Float32, reflect.Float64:
		for i := 0; i < size; i++ {
			val, _ := strconv.ParseFloat(items[i], 64)
			slice.Index(i).SetFloat(val)
		}
	case reflect.Bool:
		for i := 0; i < size; i++ {
			val, _ := strconv.ParseBool(items[i])
			slice.Index(i).SetBool(val)
		}
	default:
		failf("unsupported type of slice item: %v", t.Elem().Kind().String())
	}

	v.Set(slice)
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
		failf("unsupported type: %v", t.Kind().String())
	}
}
