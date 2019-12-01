package configuration

import (
	"os"
	"strings"
)

type valueProvider func(key string) (value interface{})

func provideFromEnv(key string) interface{} {
	v, _ := os.LookupEnv(strings.ToUpper(key))
	return v
}

var flags []interface{} //todo

func provideFromFlags(key string) interface{} {
	panic("not implemented")
	return nil
}
