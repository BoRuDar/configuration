package configuration

import (
	"os"
	"testing"
	"time"
)

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
	cfg := struct {
		Name     string `flag:"name"`
		LastName string `default:"defaultLastName"`
		Age      byte   `env:"AGE_ENV"               default:"-1"`
		BoolPtr  *bool  `default:"false"`
		ObjPtr   *struct {
			F32       float32       `default:"32"`
			StrPtr    *string       `default:"str_ptr_test"`
			HundredMS time.Duration `default:"100ms"`
		}
		Obj struct {
			IntPtr     *int16   `default:"123"`
			Beta       int      `file_json:"inside.beta"   default:"24"`
			StrSlice   []string `default:"one;two"`
			IntSlice   []int64  `default:"3; 4"`
			unexported string   `xml:"ignored"`
		}
		URLs []*string `default:"http://localhost:3000;1.2.3.4:8080"`
	}{}

	configurator := New(
		&cfg,
		// order of execution will be preserved:
		NewFlagProvider(),             // 1st
		NewEnvProvider(),              // 2nd
		NewJSONFileProvider(fileName), // 3rd
		NewDefaultProvider(),          // 4th
	)

	if err := configurator.InitValues(); err != nil {
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

	for i := range expectedURLs {
		assert(t, expectedURLs[i], *cfg.URLs[i])
	}
}

func TestConfigurator_Errors(t *testing.T) {
	tests := map[string]struct {
		input     any
		providers []Provider
	}{
		"empty providers": {
			input:     &struct{}{},
			providers: []Provider{},
		},
		"non-pointer": {
			input: struct{}{},
			providers: []Provider{
				NewDefaultProvider(),
			},
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			err := New(test.input, test.providers...).InitValues()
			if err == nil {
				t.Fatal("expected error but got nil")
			}
		})
	}
}

func TestEmbeddedFlags(t *testing.T) {
	type (
		Client struct {
			ServerAddress string `flag:"addr|127.0.0.1:443|server address"`
		}
		Config struct {
			Client *Client
		}
	)
	os.Args = []string{"smth", "-addr=addr_value"}

	var cfg Config
	if err := New(&cfg, NewFlagProvider()).InitValues(); err != nil {
		t.Fatal("unexpected err: ", err)
	}

	assert(t, true, cfg.Client != nil)
	assert(t, cfg.Client.ServerAddress, "addr_value")
}

func TestFallBackToDefault(t *testing.T) {
	// defining a struct
	cfg := struct {
		NameFlag string `flag:"name_flag||Some description"   default:"default_val"`
	}{}

	c := New(&cfg,
		NewFlagProvider(),
		NewDefaultProvider(),
	)

	if err := c.InitValues(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert(t, "default_val", cfg.NameFlag)
}

func TestSetOnFailFn(t *testing.T) {
	cfg := struct {
		Name string `default:"test_name"`
	}{}
	onFailFn := func(err error) {
		if err != nil && err.Error() != "configurator: field [Name] with tags [default:\"test_name\"] cannot be set" {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	c := New(
		&cfg,
		NewFlagProvider(),
	).
		SetOptions(
			OnFailFnOpt(onFailFn),
		)

	if err := c.InitValues(); err != nil {
		t.Fatal("unexpected err: ", err)
	}
}

func TestProviderName(t *testing.T) {
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
			assert(t, test.expectedName, test.provider.Name())
		})
	}
}

func TestConfigurator_NameCollision(t *testing.T) {
	err := New(&struct{}{}, NewDefaultProvider(), NewDefaultProvider()).InitValues()
	assert(t, ErrProviderNameCollision, err)
}

func TestConfigurator_FailedProvider(t *testing.T) {
	err := New(&struct{}{}, NewJSONFileProvider("doesn't exist")).InitValues()
	assert(t, err.Error(), "cannot init [JSONFileProvider] provider: open doesn't exist: no such file or directory")
}

func Test_FromEnvAndDefault(t *testing.T) {
	t.Setenv("AGE", "24")

	type st struct {
		Name string `env:"name"    default:"Alex"`
		Age  int    `env:"AGE"     default:"42"`
	}

	cfg := st{}

	if err := FromEnvAndDefault(&cfg); err != nil {
		t.Fatal("unexpected err", err)
	}

	assert(t, cfg.Name, "Alex")
	assert(t, cfg.Age, 24)
}
