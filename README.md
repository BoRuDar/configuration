# Configuration
is a library for injecting values recursively into structs - a convinient way of setting up a configuration object.
Currently supports next features:
- setting *default* values for struct fields (`NewDefaultProvider()`)
- setting values from *environment* variables (`NewEnvProvider()`)
- setting values from command line *flags* (`NewFlagProvider(&cfg)`)

Next fields' types are supported:
- `string`, `*string`
- `bool`, `*bool`
- `int`, `int8`, `int16`, `int32`, `int64`
- `*int`, `*int8`, `*int16`, `*int32`, `*int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `*uint`, `*uint8`, `*uint16`, `*uint32`, `*uint64`
- `float32`, `float64`,
- `*float32`, `*float64`,
- embedded structs and pointers to structs

# Quick start

```go
// define a configuration object
    cfg := struct {
        Name     string `json:"name"          default:"defaultName"         flag:"name"`
        LastName string `json:"last_name"     default:"defaultLastName"`
        Age      byte   `json:"age"           env:"AGE_ENV"`
        IsDebug  bool   `json:"is_debug"      default:"false"`
        Object   struct {
            One string  `json:"one"            default:"defaultValForOne"`
            Two float32 `json:"two"            default:"33"`
        }
        StrPtr  *string `json:"str_ptr"         default:"str_ptr_test"`
        IntPtr  *int    `json:"int_ptr"         default:"123"`
        BoolPtr *bool   `json:"bool_ptr"        default:"true"`
    }{}
    
    
    configurator, err := New(
        &cfg, // pointer to the object
        []Provider{ // list of providers
            NewFlagProvider(&cfg), // flag provider expects pointer to the object to initialize flags
            NewEnvProvider(),
            NewDefaultProvider(),
        },
        false, // logging if true
        false, // fail fast if cannot set any field in the given object
    )
    if err != nil {
        panic(err)
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
	Provide(field reflect.StructField, v reflect.Value) bool
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
If `env` tag is not found it will try to use `json` tag in upper case (`age` becomes `AGE` in this example).
Name inside tag `env:"<name>"` must be unique for each field.


### Flag provider
Looks for `flag` tag and tries to set value from the command line flag `-name`
```go
    struct {
        // ...
        Name     string `json:"name"  flag:"name"`
        // ...
    }
```
If `flag` tag is not found it will try to use value from `json` tag.
Name inside tag `flag:"<name>"` must be unique for each field.
`NewFlagProvider(&cfg)` expects a pointer to the same object for initialization.