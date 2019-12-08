package main

import (
	"flag"
	"fmt"

	"github.com/BoRuDar/configuration"
)

func main() {
	cfg := struct {
		Name     string `json:"name"  flag:"name"        default:"defaultName"`
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

	configuration.NewProviderFromFlags(&cfg)
	flag.Parse()

	err := configuration.FillUp(&cfg)
	fmt.Printf("err: %v ||| %+v", err, cfg)
}
