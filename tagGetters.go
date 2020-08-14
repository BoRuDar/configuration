package configuration

import "reflect"

func getEnvTag(f reflect.StructField) string {
	return f.Tag.Get("env")
}

func getFlagTag(f reflect.StructField) string {
	return f.Tag.Get("flag")
}

func getJSONTag(f reflect.StructField) string {
	return f.Tag.Get("json")
}

func getDefaultTag(f reflect.StructField) string {
	return f.Tag.Get("default")
}

func getValidateTag(f reflect.StructField) string {
	return f.Tag.Get("validate")
}