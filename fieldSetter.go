package configuration

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const sliceSeparator = ";"

// FieldSetter interface
type FieldSetter interface {
	SetField(field reflect.StructField, val reflect.Value, valStr string) error
}

// SetField sets field with `valStr` value (and converts it into the proper type beforehand)
func SetField(field reflect.StructField, val reflect.Value, valStr string) error {
	if val.CanInterface() {
		if fs, ok := val.Addr().Interface().(FieldSetter); ok {
			return fs.SetField(field, val, valStr) // nolint:wrapcheck
		}

		if fs, ok := val.Interface().(FieldSetter); ok {
			return fs.SetField(field, val, valStr) // nolint:wrapcheck
		}
	}

	if val.Kind() == reflect.Pointer {
		return setPtrValue(field.Type, val, valStr)
	}

	return setValue(field.Type, val, valStr)
}

func setValue(t reflect.Type, v reflect.Value, val string) error {
	var err error

	// nolint:exhaustive
	switch t.Kind() {
	case reflect.String:
		v.SetString(val)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		i, _ := strconv.ParseInt(val, 10, 64)
		v.SetInt(i)

	case reflect.Int64:
		setInt64(v, val)

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
		err = setSlice(t, v, val)

	default:
		err = fmt.Errorf("setValue: unsupported type: %v", v.Kind().String())
	}

	return err
}

func setInt64(v reflect.Value, val string) {
	// special case for parsing human-readable input for time.Duration
	if _, ok := v.Interface().(time.Duration); ok {
		d, _ := time.ParseDuration(val)
		v.SetInt(int64(d))
		return
	}

	// regular int64 case
	i, _ := strconv.ParseInt(val, 10, 64)
	v.SetInt(i)
}

// nolint:cyclop
func setSlice(t reflect.Type, v reflect.Value, val string) error {
	var (
		slice reflect.Value
		items = splitIntoSlice(val)
		size  = len(items)
	)

	if size == 0 {
		return fmt.Errorf("setSlice: got empty slice")
	}

	// nolint:exhaustive
	switch t.Elem().Kind() {
	case reflect.String:
		slice = reflect.MakeSlice(t, size, size)
		for i := range size {
			slice.Index(i).SetString(items[i])
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		slice = reflect.MakeSlice(t, size, size)
		for i := range size {
			val, _ := strconv.ParseInt(items[i], 10, 64)
			slice.Index(i).SetInt(val)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		slice = reflect.MakeSlice(t, size, size)
		for i := range size {
			val, _ := strconv.ParseUint(items[i], 10, 64)
			slice.Index(i).SetUint(val)
		}

	case reflect.Float32, reflect.Float64:
		slice = reflect.MakeSlice(t, size, size)
		for i := range size {
			val, _ := strconv.ParseFloat(items[i], 64)
			slice.Index(i).SetFloat(val)
		}
	case reflect.Bool:
		slice = reflect.MakeSlice(t, size, size)
		for i := range size {
			val, _ := strconv.ParseBool(items[i])
			slice.Index(i).SetBool(val)
		}

	case reflect.Pointer:
		slice = reflect.MakeSlice(t, size, size)
		for i := range size {
			err := setPtrValue(slice.Index(i).Type(), slice.Index(i), items[i])
			if err != nil {
				return fmt.Errorf("setSlice: cannot set type [%s] at index [%d]", slice.Index(i).Type(), i)
			}
		}

	default:
		return fmt.Errorf("setSlice: unsupported type of slice item: %v", t.Elem().Kind().String())
	}

	v.Set(slice)
	return nil
}

// nolint:cyclop
func setPtrValue(t reflect.Type, v reflect.Value, val string) error {
	var err error

	switch t.Elem().Name() {
	case reflect.Int.String():
		if i, err := strconv.Atoi(val); err == nil {
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

	case reflect.Uint.String():
		var ui uint
		if _, err := fmt.Sscan(val, &ui); err == nil {
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
		err = fmt.Errorf("setPtrValue: unsupported type: %v", t.Kind().String())
	}

	return err
}

func splitIntoSlice(val string) []string {
	var items []string

	for _, item := range strings.Split(val, sliceSeparator) {
		item = strings.TrimSpace(item)
		if len(item) > 0 {
			items = append(items, item)
		}
	}

	return items
}

func ToPtr[T any](val T) *T {
	return &val
}
