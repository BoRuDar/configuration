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
			One string `json:"one"            default:"defaultValForOne"`
		}
	}{}

	err := FillUp(&cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", cfg)
}
