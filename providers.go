package configuration

import (
	"flag"
	"log"
	"os"
	"reflect"
	"strings"
)

type valueProvider func(field reflect.StructField, v reflect.Value) bool

func provideFromDefault(field reflect.StructField, v reflect.Value) bool {
	valStr := getDefaultTag(field)
	if len(valStr) == 0 {
		return false
	}

	setField(field, v, valStr)
	return true
}

func provideFromEnv(field reflect.StructField, v reflect.Value) bool {
	key := getEnvTag(field)
	if len(key) == 0 { // if "env" is not set try to use regular json tag
		key = strings.ToUpper(getJSONTag(field))
	}
	if len(key) == 0 {
		// field doesn't have a proper tag
		return false
	}

	valStr, ok := os.LookupEnv(strings.ToUpper(key))
	if !ok || len(valStr) == 0 {
		return false
	}

	setField(field, v, valStr)
	return true
}

var flags = map[string]func() *string{}

func provideFromFlags(field reflect.StructField, v reflect.Value) bool {
	key := getFlagTag(field)
	if len(key) == 0 { // if "flag" is not set try to use regular json tag
		key = getJSONTag(field)
	}
	if len(key) == 0 {
		// field doesn't have a proper tag
		return false
	}

	if len(flags) == 0 {
		log.Println("flags: ", 0)
		return false
	}

	fn, ok := flags[key]
	if !ok {
		log.Println(key, " is not found")
		return false
	}

	val := fn()
	setField(field, v, *val)
	return true
}

func NewProviderFromFlags(i interface{}) {
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
		log.Println(t.Field(i).Tag)
		if t.Field(i).Type.Kind() == reflect.Struct {
			NewProviderFromFlags(v.Field(i).Addr().Interface())
			continue
		}
		getValFromFlags(t.Field(i))
	}
}

func getValFromFlags(field reflect.StructField) {
	log.Println("getValFromFlags: ", field.Tag)

	key := getFlagTag(field)
	if len(key) == 0 { // if "flag" is not set try to use regular json tag
		key = getJSONTag(field)
	}
	if len(key) == 0 {
		// field doesn't have a proper tag
		return
	}

	valStr := flag.String(key, "", "")

	if _, ok := flags[key]; ok {
		log.Println(key, " is already in map")
		return
	}

	flags[key] = func() *string {
		return valStr
	}
}
