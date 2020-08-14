package configuration

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
)

const flagSeparator = "|"

// NewFlagProvider creates a new provider to fetch data from flags like: --flag_name some_value
func NewFlagProvider(ptrToCfg interface{}) flagProvider {
	fp := flagProvider{
		flagsValues: map[string]func() *string{},
		flags:       map[string]*flagData{},
	}
	fp.initFlagProvider(ptrToCfg)

	flag.Parse()
	return fp
}

type flagProvider struct {
	flagsValues map[string]func() *string
	flags       map[string]*flagData
}

type flagData struct {
	key, defaultVal, usage string
}

func (fp flagProvider) initFlagProvider(i interface{}) error {
	var (
		t = reflect.TypeOf(i)
		v = reflect.ValueOf(i)
	)

	switch t.Kind() {
	case reflect.Ptr:
		t = t.Elem()
		v = v.Elem()
	default:
		return fmt.Errorf("not a pointer to a struct")
	}

	for i := 0; i < t.NumField(); i++ {
		tField := t.Field(i)
		if tField.Type.Kind() == reflect.Struct {
			if err := fp.initFlagProvider(v.Field(i).Addr().Interface()); err != nil {
				return err
			}
			continue
		}

		if tField.Type.Kind() == reflect.Ptr && tField.Type.Elem().Kind() == reflect.Struct {
			v.Field(i).Set(reflect.New(tField.Type.Elem()))

			if err := fp.initFlagProvider(v.Field(i).Interface()); err != nil {
				return err
			}
			continue
		}

		fp.setFlagCallbacks(tField)
	}
	return nil
}

func (fp flagProvider) setFlagCallbacks(field reflect.StructField) error {
	fd, err := getFlagData(field)
	if err != nil {
		return err
	}

	if _, ok := fp.flagsValues[fd.key]; ok {
		return fmt.Errorf("flagProvider: flag for the key [%s] is already set", fd.key) // TOOD: test
	}
	fp.flags[fd.key] = fd

	valStr := flag.String(fd.key, fd.defaultVal, fd.usage)
	fp.flagsValues[fd.key] = func() *string {
		return valStr
	}
	return nil
}

func (fp flagProvider) Provide(field reflect.StructField, v reflect.Value, _ ...string) error {
	fd, err := getFlagData(field)
	if err != nil {
		return err
	}

	if len(fp.flagsValues) == 0 {
		return fmt.Errorf("flagProvider: map of flagsValues is empty, nothing to fetch")
	}

	fn, ok := fp.flagsValues[fd.key]
	if !ok {
		return fmt.Errorf("flagProvider: callback for key [%s] is not found", fd.key)
	}

	val := fn()
	if *val != fd.defaultVal {
		return SetField(field, v, *val)
	}
	return fmt.Errorf("flagProvider: value for key [%s] is same as default [%s], ignoring", fd.key, fd.defaultVal)
}

func getFlagData(field reflect.StructField) (*flagData, error) {
	key := getFlagTag(field)
	if len(key) == 0 {
		return nil, fmt.Errorf("flagProvider: getFlagTag returns empty value")
	}

	flagInfo := strings.Split(key, flagSeparator)
	switch len(flagInfo) {
	case 3:
		return &flagData{
			key:        strings.TrimSpace(flagInfo[0]),
			defaultVal: strings.TrimSpace(flagInfo[1]),
			usage:      flagInfo[2],
		}, nil
	case 2:
		return &flagData{
			key:        strings.TrimSpace(flagInfo[0]),
			defaultVal: strings.TrimSpace(flagInfo[1]),
		}, nil
	case 1:
		return &flagData{
			key: strings.TrimSpace(flagInfo[0]),
		}, nil
	default:
		return nil, fmt.Errorf("flagProvider: wrong flag definition [%s]", key)
	}
}
