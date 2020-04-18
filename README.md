[![Go Report Card](https://goreportcard.com/badge/github.com/borudar/configuration)](https://goreportcard.com/report/github.com/borudar/configuration)
[![codecov](https://codecov.io/gh/BoRuDar/configuration/branch/master/graph/badge.svg)](https://codecov.io/gh/BoRuDar/configuration)
![Go](https://github.com/BoRuDar/configuration/workflows/Go/badge.svg)
[![GoDoc](https://godoc.org/github.com/BoRuDar/configuration?status.png)](https://godoc.org/github.com/BoRuDar/configuration)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go) 

# Configuration
is a library for injecting values recursively into structs - a convenient way of setting up a configuration object.
Available features:
- setting *default* values for struct fields - `NewDefaultProvider()`
- setting values from *environment* variables - `NewEnvProvider()`
- setting values from command line *flags* - `NewFlagProvider(&cfg)`
- setting values from *files* (JSON or YAML) - `NewFileProvider("./testdata/input.yml")`

Supported types:
- `string`, `*string`, `[]string`
- `bool`, `*bool`, `[]bool`
- `int`, `int8`, `int16`, `int32`, `int64` + slices of these types
- `*int`, `*int8`, `*int16`, `*int32`, `*int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64` + slices of these types
- `*uint`, `*uint8`, `*uint16`, `*uint32`, `*uint64`
- `float32`, `float64` + slices of these types
- `*float32`, `*float64`
- `time.Duration` from strings like `12ms`, `2s` etc.
- embedded structs and pointers to structs

# Quick start
Import path `github.com/BoRuDar/configuration/v2`
```go
// define a configuration object
cfg := struct {
    Name     string `json:"name"          default:"defaultName"         flag:"name"`
    LastName string `json:"last_name"     default:"defaultLastName"`
    Age      byte   `json:"age"           env:"AGE_ENV"`
    BoolPtr  *bool  `json:"bool_ptr"      default:"false"`

    ObjPtr *struct {
        F32       float32       `json:"f32"            default:"32"`
        StrPtr    *string       `json:"str_ptr"        default:"str_ptr_test"`
        HundredMS time.Duration `json:"hundred_ms"     default:"100ms"`
    }

    Obj struct {
        IntPtr   *int16   `json:"int_ptr"         default:"123"`
        NameYML  int      `default:"24"`
        StrSlice []string `default:"one;two"`
        IntSlice []int64  `default:"3; 4"`
    }
}{}

configurator, err := New(
    &cfg, // pointer to the object for configuration 
    NewFlagProvider(&cfg),  // 1. flag provider expects pointer to the object to initialize flags
    NewEnvProvider(),       // 2.
    NewFileProvider("./testdata/input.yml"), // 3.
    NewDefaultProvider(),   // 4.
    // providers are executed in order of the declaration from 1 to 4 
)
if err != nil {
    t.Fatalf("unexpected error: %v", err)
}

configurator.InitValues()
```


# Providers
You can specify one or more providers. They will be executed in order of definition:
```go
[]Provider{
    NewFlagProvider(&cfg), // 1
    NewEnvProvider(), // 2
    NewDefaultProvider(), // 3
} 
```
If provider set value successfully next ones will not be executed (if flag provider from the sample above found a value env and default providers are skipped). 
The value of first successfully executed provider will be set.
If none of providers found value - "zero" value of a field remains.
You can define a custom provider which should satisfy next interface:
```go
type Provider interface {
	Provide(field reflect.StructField, v reflect.Value, pathToField ...string) bool
}
```

### Default provider
Looks for `default` tag and set value from it:
```go
    struct {
        // ...
        Name string `json:"name"          default:"defaultName"`
        // ...
    }
```


### Env provider
Looks for `env` tag and tries to find an ENV variable with the name from the tag (`AGE_ENV` in this example):
```go
    struct {
        // ...
        Age      byte   `json:"age"           env:"AGE_ENV"`
        // ...
    }
```
Name inside tag `env:"<name>"` must be unique for each field.


### Flag provider
Looks for `flag` tag and tries to set value from the command line flag `-name`
```go
    struct {
        // ...
        Name     string `json:"name"  flag:"name|default_value|Description"`
        // ...
    }
```
Name inside tag `flag:"<name>"` must be unique for each field. `default_value` and `description` sections are optional and can be omitted.
`NewFlagProvider(&cfg)` expects a pointer to the same object for initialization.

*Note*: if program is executed with `-help` or `-h` flag you will see all available flags with description:
```bash
Flags: 
	-flag_name		"Description (default: default_value)"
``` 
And program execution will be terminated.

### File provider
Doesn't require any specific tags. JSON and YAML formats of files are supported.
```go
    NewFileProvider("./testdata/input.yml")
```
