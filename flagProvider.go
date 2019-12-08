package configuration

import (
	"flag"
	"reflect"
)

func NewFlagProvider(i interface{}) flagProvider {
	fp := flagProvider{flags: map[string]func() *string{}}
	fp.initFlagProvider(i)
	flag.Parse()
	return fp
}

type flagProvider struct {
	flags map[string]func() *string
}

func (fp flagProvider) initFlagProvider(i interface{}) {
	var (
		t = reflect.TypeOf(i)
		v = reflect.ValueOf(i)
	)

	switch t.Kind() {
	case reflect.Ptr:
		t = t.Elem()
		v = v.Elem()
	default:
		panic("not a pointer to a struct")
	}

	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Type.Kind() == reflect.Struct {
			fp.initFlagProvider(v.Field(i).Addr().Interface())
			continue
		}
		fp.getValFromFlags(t.Field(i))
	}
}

func (fp flagProvider) getValFromFlags(field reflect.StructField) {
	key := getFlagTag(field)
	if len(key) == 0 { // if "flag" is not set try to use regular json tag
		key = getJSONTag(field)
	}
	if len(key) == 0 {
		// field doesn't have a proper tag
		return
	}

	if _, ok := fp.flags[key]; ok {
		return
	}

	valStr := flag.String(key, "", "")
	fp.flags[key] = func() *string {
		return valStr
	}
}

func (fp flagProvider) Provide(field reflect.StructField, v reflect.Value) bool {
	key := getFlagTag(field)
	if len(key) == 0 { // if "flag" is not set try to use regular json tag
		key = getJSONTag(field)
	}
	if len(key) == 0 {
		// field doesn't have a proper tag
		return false
	}

	if len(fp.flags) == 0 {
		return false
	}

	fn, ok := fp.flags[key]
	if !ok {
		return false
	}

	val := fn()
	setField(field, v, *val)
	return len(*val) > 0
}
