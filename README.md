[![Go Report Card](https://goreportcard.com/badge/github.com/BoRuDar/configuration/v4)](https://goreportcard.com/report/github.com/BoRuDar/configuration/v4)
[![codecov](https://codecov.io/gh/BoRuDar/configuration/branch/master/graph/badge.svg)](https://codecov.io/gh/BoRuDar/configuration)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/6295/badge)](https://www.bestpractices.dev/projects/6295)
[![GoDoc](https://godoc.org/github.com/BoRuDar/configuration?status.png)](https://godoc.org/github.com/BoRuDar/configuration/v4)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)


# Configuration
is a library for injecting values recursively into structs - a convenient way of setting up a configuration object.
Available features:
- setting *default* values for struct fields - `NewDefaultProvider()`
- setting values from *environment* variables - `NewEnvProvider()`
- setting values from command line *flags* - `NewFlagProvider(&cfg)`
- setting values from a JSON *file* - `NewJSONFileProvider("./testdata/input.json")`

## Supported types:
- `string`, `*string`, `[]string`, `[]*string`
- `bool`, `*bool`, `[]bool`, `[]*bool`
- `int`, `int8`, `int16`, `int32`, `int64` + slices of these types
- `*int`, `*int8`, `*int16`, `*int32`, `*int64` + slices of these types
- `uint`, `uint8`, `uint16`, `uint32`, `uint64` + slices of these types
- `*uint`, `*uint8`, `*uint16`, `*uint32`, `*uint64` + slices of these types
- `float32`, `float64` + slices of these types
- `*float32`, `*float64` + slices of these types
- `time.Duration` from strings like `12ms`, `2s` etc.
- embedded structs and pointers to structs
- any custom type which satisfies `FieldSetter` [interface](#FieldSetter-interface)


# Why?
- your entire configuration can be defined in one model
- all metadata is in your model (defined with `tags`)
- easy to set/change a source of data for your configuration
- easy to set a priority of sources to fetch data from (e.g., 1.`flags`, 2.`env`, 3.`default` or another order)
- you can implement your custom provider
- no external dependencies
- complies with `12-factor app`


# Quick start
Import path `github.com/BoRuDar/configuration/v5`
```go
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

cfg, err := New[Conf](  // specify the [T] of the structure to be returned
    // order of execution will be preserved:
    NewFlagProvider(),             // 1st
    NewEnvProvider(),              // 2nd
    NewJSONFileProvider(fileName), // 3rd
    NewDefaultProvider(),          // 4th
)
if err != nil {
    t.Fatalf("unexpected error: %v", err)
}
```

If you need only ENV variables and default values you can use a shorter form:
```go
cfg, err := configuration.FromEnvAndDefault[T]()
```


# Providers
You can specify one or more providers. They will be executed in order of definition:
```go
[]Provider{
    NewFlagProvider(),     // 1
    NewEnvProvider(),      // 2
    NewDefaultProvider(),  // 3
} 
```
**IMPORTANT:** If provider sets value successfully next ones will **NOT** be executed 
(if flag provider from the sample above finds the value - then the env and default providers are skipped). 
The value of the first successfully executed provider will be set.
If none of providers can set value - an error will be returned.


### Custom provider
You can define a custom provider which should satisfy this interface:
```go
type Provider interface {
    Name() string
    Tag() string
    Init(ptr any) error
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
So `Name` will be set to "defaultName".


### Env provider
Looks for `env` tag and tries to find an ENV variable with the name from the tag (`AGE` in this example):
```go
struct {
    // ...
    Age      byte   `env:"AGE"`
    // ...
}
```
Name inside tag `env:"<name>"` must be unique for each field. 
Only strings in **UPPER** register for ENV vars are accepted:
```bash
bad_env_var_name=bad
Also_Bad_Env_Var_Name=bad
GOOD_ENV_VAR_NAME=good
```


### Flag provider
Looks for `flag` tag and tries to set the value from the command line flag `-first_name`
```go
struct {
    // ...
    Name     string `flag:"first_name|default_value|Description"`
    // ...
}
```
Name inside tag `flag:"<name>"` must be unique for each field.
`default_value` and `description` sections are `optional` and may be omitted.

*Note*: if program is executed with `-help` or `-h` flag you will see all available flags with description:
```bash
Flags: 
	-first_name		"Description (default: default_value)"
``` 
And program execution will be terminated.
#### Options for _NewFlagProvider_
* `WithFlagSet(s FlagSet)`  - sets a custom `FlagSet`


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

### Additional providers
* [YAML files](https://github.com/BoRuDar/configuration-yaml-file)


## FieldSetter interface
You can define how to set fields with any custom types: 
```go
type FieldSetter interface {
	SetField(field reflect.StructField, val reflect.Value, valStr string) error
}
```
Example:
```go
type ipTest net.IP

func (it *ipTest) SetField(_ reflect.StructField, val reflect.Value, valStr string) error {
	i := ipTest(net.ParseIP(valStr))

	if val.Kind() == reflect.Pointer {
		val.Set(reflect.ValueOf(&i))
	} else {
		val.Set(reflect.ValueOf(i))
	}

	return nil
}
```


# Contribution
1. Open a feature request or a bug report in [issues](https://github.com/BoRuDar/configuration/issues)
2. Fork and create a PR into `dev` branch
