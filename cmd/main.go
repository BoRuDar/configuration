package main

import (
	"fmt"

	"github.com/BoRuDar/configuration"
)

func main() {
	cfg := struct {
		Name     string `json:"name"          default:"defaultName"         flag:"name"`
		LastName string `json:"last_name"     default:"defaultLastName"`
		Age      byte   `json:"age"           env:"AGE_ENV"`
		IsDebug  bool   `json:"is_debug"`
		Obj      struct {
			One string  `json:"one"            default:"defaultValForOne"`
			Two float32 `json:"two"            default:"33"`
		}
		StrPtr  *string `json:"str_ptr"         default:"str_ptr_test"`
		IntPtr  *int    `json:"int_ptr"         default:"123"`
		BoolPtr *bool   `json:"bool_ptr"        default:"true"`
	}{}

	configurator, err := configuration.New(&cfg, []configuration.Provider{
		configuration.NewFlagProvider(&cfg),
		configuration.NewEnvProvider(),
		configuration.NewDefaultProvider(),
	})
	if err != nil {
		panic(err)
	}
	if err = configurator.FillUp(); err != nil {
		panic(err)
	}

	fmt.Printf("%+v", cfg)
}
