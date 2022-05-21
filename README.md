[![Go Report Card](https://goreportcard.com/badge/github.com/borudar/configuration)](https://goreportcard.com/report/github.com/borudar/configuration)
[![codecov](https://codecov.io/gh/BoRuDar/configuration/branch/master/graph/badge.svg)](https://codecov.io/gh/BoRuDar/configuration)
![Go](https://github.com/BoRuDar/configuration/workflows/Go/badge.svg)
[![GoDoc](https://godoc.org/github.com/BoRuDar/configuration?status.png)](https://godoc.org/github.com/BoRuDar/configuration/v3)
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

# Why?
- your entire configuration can be defined in one model
- all metadata is in your model (defined with `tags`)
- easy to set/change a source of data for your configuration
- easy to set a priority of sources to fetch data from (e.g., 1.`flags`, 2.`env`, 3.`default` or another order)
- you can implement your custom provider
- no external dependencies
- complies with `12-factor app`

# Quick start
Import path `github.com/BoRuDar/configuration/v3`
```go
	// defining a struct
cfg := struct {
    Name     string `flag:"name"`
    LastName string `default:"defaultLastName"`
    Age      byte   `env:"AGE_ENV"    default:"-1"`
    BoolPtr  *bool  `default:"false"`
    
    ObjPtr *struct {
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
}{}

configurator := configuration.New(
    &cfg,
    // order of execution will be preserved: 
    configuration.NewFlagProvider(),             // 1st
    configuration.NewEnvProvider(),              // 2nd 
    configuration.NewJSONFileProvider(fileName), // 3rd 
    configuration.NewDefaultProvider(),          // 4th
)

if err := configurator.InitValues(); err != nil {
    log.Fatalf("unexpected error: ", err)
}
```

If you need only ENV variables and default values you can use a shorter form:
```go
err := configuration.FromEnvAndDefault(&cfg)
```


# Providers
You can specify one or more providers. They will be executed in order of definition:
```go
[]Provider{
    NewFlagProvider(&cfg), // 1
    NewEnvProvider(),      // 2
    NewDefaultProvider(),  // 3
} 
```
If provider set value successfully next ones will not be executed (if flag provider from the sample above found a value env and default providers are skipped). 
The value of first successfully executed provider will be set.
If none of providers found value - an application will be terminated.
This behavior can be changed with `configurator.OnFailFnOpt` option:
```go
err := configuration.New(
    &cfg,
    configuration.NewEnvProvider(),
    configuration.NewDefaultProvider()).
    SetOptions(
        configuration.OnFailFnOpt(func(err error) {
            log.Println(err)
        }),
    ).InitValues()
```


### Custom provider
You can define a custom provider which should satisfy next interface:
```go
type Provider interface {
    Name() string
    Init(ptr interface{}) error
    Provide(field reflect.StructField, v reflect.Value) error
}
```

### Default provider
Looks for `default` tag and set value from it:
```go
struct {
    // ...
    Name string `default:"defaultName"`
    // ...
}
```


### Env provider
Looks for `env` tag and tries to find an ENV variable with the name from the tag (`AGE` in this example):
```go
struct {
    // ...
    Age      byte   `env:"AGE"`
    // ...
}
```
Name inside tag `env:"<name>"` must be unique for each field. Only UPPER register for ENV vars is accepted:
```bash
bad_env_var_name=bad
GOOD_ENV_VAR_NAME=good
```


### Flag provider
Looks for `flag` tag and tries to set value from the command line flag `-first_name`
```go
struct {
    // ...
    Name     string `flag:"first_name|default_value|Description"`
    // ...
}
```
Name inside tag `flag:"<name>"` must be unique for each field. `default_value` and `description` sections are `optional` and can be omitted.
`NewFlagProvider(&cfg)` expects a pointer to the same object for initialization.

*Note*: if program is executed with `-help` or `-h` flag you will see all available flags with description:
```bash
Flags: 
	-first_name		"Description (default: default_value)"
``` 
And program execution will be terminated.
#### Options for _NewFlagProvider_
* WithFlagSet - sets a custom FlagSet


### JSON File provider 
Requires `file_json:"<path_to_json_field>"` tag.
```go
NewJSONFileProvider("./testdata/input.json")
```
For example, tag `file_json:"cache.retention"` will assume that you have this structure of your JSON file:
```json
{
  "cache": {
    "retention": 1
  }
}
```
