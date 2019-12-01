package configuration

import (
	"testing"
)

func TestName(t *testing.T) {
	cfg := struct {
		Name     string `json:"name"          default:"test"`
		LastName string `json:"last_name"     default:"defaultVal"`
		Age      int16  `json:"age"`
		IsDebug  bool   `json:"is_debug"`
		Obj      struct {
			One string  `json:"one"            default:"defaultValForOne"`
			Two float32 `json:"two"            default:"33"`
		}
		StrPtr  *string `json:"str_ptr"         default:"str_ptr_test"`
		IntPtr  *int    `json:"int_ptr"         default:"234"`
		BoolPtr *bool   `json:"bool_ptr"        default:"true"`
	}{}

	err := FillUp(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", cfg)
	t.Log(*cfg.StrPtr)
	t.Log(*cfg.IntPtr)
	t.Log(cfg.BoolPtr)
}
