package configuration

import (
	"flag"
	"reflect"
)

func NewFlagProvider(ptrToCfg interface{}) flagProvider {
	fp := flagProvider{flags: map[string]func() *string{}}
	fp.initFlagProvider(ptrToCfg)
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
	if len(key) == 0 {
		logf("flagProvider: getFlagTag returns empty value")
		// if "flag" is not set try to use regular json tag
		key = getJSONTag(field)
	}
	if len(key) == 0 {
		logf("flagProvider: key [%s] is empty", key)
		// field doesn't have a proper tag
		return
	}

	if _, ok := fp.flags[key]; ok {
		logf("flagProvider: cannot find value for key [%s]", key)
		return
	}

	valStr := flag.String(key, "", "")
	fp.flags[key] = func() *string {
		return valStr
	}
}

func (fp flagProvider) Provide(field reflect.StructField, v reflect.Value) bool {
	key := getFlagTag(field)
	if len(key) == 0 {
		logf("flagProvider: getFlagTag returns empty value")
		// if "flag" is not set try to use regular json tag
		key = getJSONTag(field)
	}
	if len(key) == 0 {
		logf("flagProvider: key is empty")
		// field doesn't have proper tags
		return false
	}

	if len(fp.flags) == 0 {
		logf("flagProvider: map of flags is empty, nothing to fetch")
		return false
	}

	fn, ok := fp.flags[key]
	if !ok {
		logf("flagProvider: flag for key [%s] already exists", key)
		return false
	}

	val := fn()
	setField(field, v, *val)
	logf("flagProvider: set [%s] to field [%s] with tags [%v]", *val, field.Name, field.Tag)
	return len(*val) > 0
}
