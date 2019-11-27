package configuration

import "reflect"

func getJSONTag(f reflect.StructField) string {
	return f.Tag.Get("json")
}

func getDefaultTag(f reflect.StructField) string {
	return f.Tag.Get("default")
}
