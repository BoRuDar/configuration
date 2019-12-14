package configuration

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigurator(t *testing.T) {
	// setting command line flag
	os.Args = []string{"smth", "-name=flag_value"}

	// setting env variable
	removeEnvKey, err := setEnv("AGE_ENV", "45")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer removeEnvKey()

	// defining a struct
	cfg := struct {
		Name     string `json:"name"          default:"defaultName"         flag:"name"`
		LastName string `json:"last_name"     default:"defaultLastName"`
		Age      byte   `json:"age"           env:"AGE_ENV"`
		BoolPtr  *bool  `json:"bool_ptr"      default:"false"`

		ObjPtr *struct {
			F32    float32 `json:"f32"            default:"32"`
			StrPtr *string `json:"str_ptr"        default:"str_ptr_test"`
		}

		Obj struct {
			IntPtr *int16 `json:"int_ptr"         default:"123"`
		}
	}{}

	configurator, err := New(&cfg, []Provider{
		NewFlagProvider(&cfg),
		NewEnvProvider(),
		NewDefaultProvider(),
	}, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	configurator.InitValues()

	assert.Equal(t, "flag_value", cfg.Name)
	assert.Equal(t, "defaultLastName", cfg.LastName)
	assert.Equal(t, byte(45), cfg.Age)
	assert.NotNil(t, cfg.BoolPtr)
	assert.Equal(t, false, *cfg.BoolPtr)

	assert.NotNil(t, cfg.ObjPtr)
	assert.Equal(t, float32(32), cfg.ObjPtr.F32)
	assert.NotNil(t, cfg.ObjPtr.StrPtr)
	assert.Equal(t, "str_ptr_test", *cfg.ObjPtr.StrPtr)

	assert.NotNil(t, cfg.Obj.IntPtr)
	assert.Equal(t, int16(123), *cfg.Obj.IntPtr)
}
