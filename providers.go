package configuration

import (
	"flag"
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

var flags []interface{} //todo

func provideFromFlags(field reflect.StructField, v reflect.Value) bool {
	key := getFlagTag(field)
	if len(key) == 0 { // if "flag" is not set try to use regular json tag
		key = getJSONTag(field)
	}
	if len(key) == 0 {
		// field doesn't have a proper tag
		return false
	}

	valStr := flag.String(key, "", "")
	if valStr == nil || len(*valStr) == 0 {
		return false
	}

	setField(field, v, *valStr)
	return true
}
