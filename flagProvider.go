package configuration

import (
	"flag"
	"reflect"
	"strings"
)

const flagSeparator = "|"

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

func (fp flagProvider) initFlagProvider(i interface{}) {
	var (
		t = reflect.TypeOf(i)
		v = reflect.ValueOf(i)
	)

	switch t.Kind() {
	case reflect.Ptr:
		t = t.Elem()
		v = v.Elem()
	default:
		failf("not a pointer to a struct")
		return
	}

	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Type.Kind() == reflect.Struct {
			fp.initFlagProvider(v.Field(i).Addr().Interface())
			continue
		}
		fp.setFlagCallbacks(t.Field(i))
	}
}

func (fp flagProvider) setFlagCallbacks(field reflect.StructField) {
	fd := getFlagData(field)
	if fd == nil {
		return
	}

	if _, ok := fp.flagsValues[fd.key]; ok {
		logf("flagProvider: flag for the key [%s] is already set", fd.key)
		return
	}
	fp.flags[fd.key] = fd

	valStr := flag.String(fd.key, fd.defaultVal, fd.usage)
	fp.flagsValues[fd.key] = func() *string {
		return valStr
	}
}

func (fp flagProvider) Provide(field reflect.StructField, v reflect.Value, _ ...string) bool {
	fd := getFlagData(field)
	if fd == nil {
		return false
	}

	if len(fp.flagsValues) == 0 {
		logf("flagProvider: map of flagsValues is empty, nothing to fetch")
		return false
	}

	fn, ok := fp.flagsValues[fd.key]
	if !ok {
		logf("flagProvider: callback for key [%s] is not found", fd.key)
		return false
	}

	val := fn()
	SetField(field, v, *val)
	logf("flagProvider: set [%s] to field [%s] with tags [%v]", *val, field.Name, field.Tag)
	return len(*val) > 0
}

func getFlagData(field reflect.StructField) *flagData {
	key := getFlagTag(field)
	if len(key) == 0 {
		logf("flagProvider: getFlagTag returns empty value")
		return nil
	}

	flagInfo := strings.Split(key, flagSeparator)
	switch len(flagInfo) {
	case 3:
		return &flagData{
			key:        strings.TrimSpace(flagInfo[0]),
			defaultVal: strings.TrimSpace(flagInfo[1]),
			usage:      flagInfo[2],
		}
	case 2:
		return &flagData{
			key:        strings.TrimSpace(flagInfo[0]),
			defaultVal: strings.TrimSpace(flagInfo[1]),
		}
	case 1:
		return &flagData{
			key: strings.TrimSpace(flagInfo[0]),
		}
	default:
		logf("flagProvider: wrong flag definition [%s]", key)
		return nil
	}
}
