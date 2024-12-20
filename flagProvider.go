package configuration

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
)

const (
	FlagProviderName = `FlagProvider`
	FlagProviderTag  = `flag`
	flagSeparator    = "|"
)

type FlagProviderOption func(*flagProvider)

// NewFlagProvider creates a new provider to fetch data from flags like: --flag_name some_value
// nolint:revive
func NewFlagProvider(opts ...FlagProviderOption) flagProvider {
	fp := flagProvider{
		flagsValues: map[string]func() *string{},
		flags:       map[string]*flagData{},
		flagSet:     flag.CommandLine,
	}

	for _, f := range opts {
		f(&fp)
	}

	return fp
}

func (flagProvider) Name() string {
	return FlagProviderName
}

func (flagProvider) Tag() string {
	return FlagProviderTag
}

func (fp flagProvider) Init(ptr any) (err error) {
	if err := fp.initFlagProvider(ptr); err != nil {
		return err
	}

	if err := fp.flagSet.Parse(os.Args[1:]); err != nil {
		return fmt.Errorf("%s.Init: %w", FlagProviderName, err)
	}

	return nil
}

// FlagSet is the part of flag.FlagSet that NewFlagProvider uses
type FlagSet interface {
	Parse(arguments []string) error
	String(name string, value string, usage string) *string
}

// WithFlagSet allows the flag.FlagSet to be provided to NewFlagProvider.
// This allows compatibility with other flag parsing utilities.
func WithFlagSet(s FlagSet) FlagProviderOption {
	return func(fp *flagProvider) {
		fp.flagSet = s
	}
}

type flagProvider struct {
	flagsValues map[string]func() *string
	flags       map[string]*flagData
	flagSet     FlagSet
}

type flagData struct {
	key        string
	defaultVal string
	usage      string
}

func (fp flagProvider) initFlagProvider(ptr any) error {
	var (
		t = reflect.TypeOf(ptr)
		v = reflect.ValueOf(ptr)
	)

	// nolint:exhaustive
	switch t.Kind() {
	case reflect.Ptr:
		t = t.Elem()
		v = v.Elem()
	default:
		return ErrInvalidInput
	}

	for i := range t.NumField() {
		tField := t.Field(i)
		if tField.Type.Kind() == reflect.Struct {
			_ = fp.initFlagProvider(v.Field(i).Addr().Interface())
			continue
		}

		if tField.Type.Kind() == reflect.Ptr && tField.Type.Elem().Kind() == reflect.Struct {
			v.Field(i).Set(reflect.New(tField.Type.Elem()))

			_ = fp.initFlagProvider(v.Field(i).Interface())
			continue
		}

		if err := fp.setFlagCallbacks(tField); err != nil && !errors.Is(err, ErrNoTag) { // 'flag' tag is not set for struct field
			return err
		}
	}
	return nil
}

func (fp flagProvider) setFlagCallbacks(field reflect.StructField) error {
	fd, err := fp.getFlagData(field)
	if err != nil {
		return err
	}

	if _, ok := fp.flagsValues[fd.key]; ok {
		return fmt.Errorf("%w: %s", ErrTagNotUnique, fd.key)
	}
	fp.flags[fd.key] = fd

	valStr := fp.flagSet.String(fd.key, fd.defaultVal, fd.usage)
	fp.flagsValues[fd.key] = func() *string {
		return valStr
	}

	return nil
}

func (fp flagProvider) Provide(field reflect.StructField, v reflect.Value) error {
	fd, err := fp.getFlagData(field)
	if err != nil {
		return err
	}

	fn := fp.flagsValues[fd.key]

	val := fn()
	if val == nil || len(*val) == 0 {
		return ErrEmptyValue
	}

	return SetField(field, v, *val)
}

func (fp flagProvider) getFlagData(field reflect.StructField) (*flagData, error) {
	key := field.Tag.Get(FlagProviderTag)
	if len(key) == 0 {
		return nil, ErrNoTag
	}

	flagInfo := strings.Split(key, flagSeparator)
	switch len(flagInfo) {
	case 3: // nolint:mnd
		return &flagData{
			key:        strings.TrimSpace(flagInfo[0]),
			defaultVal: strings.TrimSpace(flagInfo[1]),
			usage:      flagInfo[2],
		}, nil

	case 2: // nolint:mnd
		return &flagData{
			key:        strings.TrimSpace(flagInfo[0]),
			defaultVal: strings.TrimSpace(flagInfo[1]),
		}, nil

	case 1: // nolint:mnd
		return &flagData{
			key: strings.TrimSpace(flagInfo[0]),
		}, nil

	default:
		return nil, fmt.Errorf("wrong flag definition [%s]", key)
	}
}
