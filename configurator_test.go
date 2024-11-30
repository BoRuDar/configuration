package configuration

import (
	"net"
	"os"
	"reflect"
	"testing"
	"time"
)

// nolint:paralleltest
// t.Setenv doesn't work with t.Parallel()
func TestConfigurator(t *testing.T) {
	// setting command line flag
	os.Args = []string{"smth", "-name=flag_value"}

	// test file
	fileName := "./testdata/input.json"

	// setting env variable
	t.Setenv("AGE_ENV", "45")

	expectedURLs := []string{
		"http://localhost:3000",
		"1.2.3.4:8080",
	}

	// defining a struct
	type Conf struct {
		Name     string `flag:"name"`
		LastName string `default:"defaultLastName"`
		Age      byte   `env:"AGE_ENV"               default:"-1"`
		BoolPtr  *bool  `default:"false"`
		ObjPtr   *struct {
			F32       float32       `default:"32"`
			StrPtr    *string       `default:"str_ptr_test"`
			HundredMS time.Duration `default:"100ms"` // nolint:stylecheck
		}
		Obj struct {
			IntPtr     *int16   `default:"123"`
			Beta       int      `file_json:"inside.beta"   default:"24"`
			StrSlice   []string `default:"one;two"`
			IntSlice   []int64  `default:"3; 4"`
			unexported string   // ignored
		}
		URLs   []*string `default:"http://localhost:3000;1.2.3.4:8080"`
		HostIP ipTest    `default:"127.0.0.3"`
	}

	configurator := New[Conf](
		// order of execution will be preserved:
		NewFlagProvider(),             // 1st
		NewEnvProvider(),              // 2nd
		NewJSONFileProvider(fileName), // 3rd
		NewDefaultProvider(),          // 4th
	)

	cfg, err := configurator.InitValues()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert(t, "flag_value", cfg.Name)
	assert(t, "defaultLastName", cfg.LastName)
	assert(t, byte(45), cfg.Age)
	assert(t, true, cfg.BoolPtr != nil)
	assert(t, false, *cfg.BoolPtr)

	assert(t, true, cfg.ObjPtr != nil)
	assert(t, float32(32), cfg.ObjPtr.F32)
	assert(t, true, cfg.ObjPtr.StrPtr != nil)
	assert(t, "str_ptr_test", *cfg.ObjPtr.StrPtr)

	assert(t, true, cfg.Obj.IntPtr != nil)
	assert(t, int16(123), *cfg.Obj.IntPtr)
	assert(t, int(42), cfg.Obj.Beta)
	assert(t, []string{"one", "two"}, cfg.Obj.StrSlice)
	assert(t, []int64{3, 4}, cfg.Obj.IntSlice)
	assert(t, time.Millisecond*100, cfg.ObjPtr.HundredMS)

	assert(t, net.ParseIP("127.0.0.3"), net.IP(cfg.HostIP))

	for i := range expectedURLs {
		assert(t, expectedURLs[i], *cfg.URLs[i])
	}
}

//func TestConfigurator_Errors(t *testing.T) {
//	t.Parallel()
//
//	tests := map[string]struct {
//		input     any
//		providers []Provider
//	}{
//		"empty providers": {
//			input:     &struct{}{},
//			providers: []Provider{},
//		},
//		"non-pointer": {
//			input: struct{}{},
//			providers: []Provider{
//				NewDefaultProvider(),
//			},
//		},
//	}
//
//	for name, test := range tests {
//		test := test
//		t.Run(name, func(t *testing.T) {
//			t.Parallel()
//
//			_, err := New[any](test.input, test.providers...).InitValues()
//			if err == nil {
//				t.Fatal("expected error but got nil")
//			}
//		})
//	}
//}

func TestEmbeddedFlags(t *testing.T) {
	t.Parallel()

	type (
		Client struct {
			ServerAddress string `flag:"addr|127.0.0.1:443|server address"`
		}
		Config struct {
			Client *Client
		}
	)
	os.Args = []string{"smth", "-addr=addr_value"}

	cfg, err := New[Config](NewFlagProvider()).InitValues()
	if err != nil {
		t.Fatal("unexpected err: ", err)
	}

	assert(t, true, cfg.Client != nil)
	assert(t, cfg.Client.ServerAddress, "addr_value")
}

// nolint:paralleltest
func TestFallBackToDefault(t *testing.T) {
	// defining a struct
	type Cfg struct {
		NameFlag string `flag:"name_flag||Some description" default:"default_val"`
	}

	c := New[Cfg](
		NewFlagProvider(),
		NewDefaultProvider(),
	)

	cfg, err := c.InitValues()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert(t, "default_val", cfg.NameFlag)
}

func TestSetOnFailFn(t *testing.T) {
	t.Parallel()

	type Cfg struct {
		Name string `flag:""`
	}
	onFailFn := func(field reflect.StructField, err error) {
		if err != nil && err.Error() != `no tag` {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	c := New[Cfg](
		NewFlagProvider(),
	).SetOptions(
		OnFailFnOpt[Cfg](onFailFn),
	)

	if _, err := c.InitValues(); err != nil {
		t.Fatal("unexpected err: ", err)
	}
}

func TestProviderName(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		provider     Provider
		expectedName string
	}{
		DefaultProviderName: {
			provider:     NewDefaultProvider(),
			expectedName: DefaultProviderName,
		},
		EnvProviderName: {
			provider:     NewEnvProvider(),
			expectedName: EnvProviderName,
		},
		FlagProviderName: {
			provider:     NewFlagProvider(),
			expectedName: FlagProviderName,
		},
		JSONFileProviderName: {
			provider:     NewJSONFileProvider(""),
			expectedName: JSONFileProviderName,
		},
	}

	for name, test := range testCases {
		test := test

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert(t, test.expectedName, test.provider.Name())
		})
	}
}

func TestConfigurator_NameCollision(t *testing.T) {
	t.Parallel()

	_, err := New[struct{}](NewDefaultProvider(), NewDefaultProvider()).InitValues()
	assert(t, ErrProviderNameCollision, err)
}

func TestConfigurator_FailedProvider(t *testing.T) {
	t.Parallel()

	_, err := New[struct{}](NewJSONFileProvider("doesn't exist")).InitValues()
	assert(t, "cannot init [JSONFileProvider] provider: JSONFileProvider.Init: open doesn't exist: no such file or directory", err.Error())
}

// nolint:paralleltest
// t.Setenv doesn't work with t.Parallel()
func Test_FromEnvAndDefault(t *testing.T) {
	t.Setenv("AGE", "24")

	type st struct {
		Name string `env:"name"    default:"Alex"`
		Age  int    `env:"AGE"     default:"42"`
	}

	cfg, err := FromEnvAndDefault[st]()
	if err != nil {
		t.Fatal("unexpected err", err)
	}

	assert(t, cfg.Name, "Alex")
	assert(t, cfg.Age, 24)
}
