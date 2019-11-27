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

func provideFromFlags(key string) interface{} {
	return nil
}
