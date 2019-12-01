package configuration

import "reflect"

func getFieldTag(f reflect.StructField) string {
	panic("not implemented")
	return "" //todo
}

func getEnvTag(f reflect.StructField) string {
	return f.Tag.Get("env")
}

func getJSONTag(f reflect.StructField) string {
	return f.Tag.Get("json")
}

func getDefaultTag(f reflect.StructField) string {
	return f.Tag.Get("default")
}
