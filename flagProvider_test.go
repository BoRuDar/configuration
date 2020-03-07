package configuration

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlagProvider(t *testing.T) {
	type testStruct struct {
		Name string `flag:"flag_name"`
	}
	testObj := testStruct{}
	os.Args = []string{"smth", "-flag_name=flag_value"}

	fieldType := reflect.TypeOf(&testObj).Elem().Field(0)
	fieldVal := reflect.ValueOf(&testObj).Elem().Field(0)

	provider := NewFlagProvider(&testObj)
	testValue := "flag_value"

	if !provider.Provide(fieldType, fieldVal) {
		t.Fatal("cannot set value")
	}

	assert.Equal(t, testValue, testObj.Name)
	assert.NotNil(t, showHelpFlag)
}

func TestFlagData_String(t *testing.T) {
	tests := map[string]struct {
		flag     *flagData
		expected string
	}{
		"3 fields": {
			flag: &flagData{
				key:        "key",
				defaultVal: "defVal",
				usage:      "usage",
			},
			expected: `	-key		"usage (default: defVal)"`,
		},
		"2 fields": {
			flag: &flagData{
				key:        "key",
				defaultVal: "defVal",
			},
			expected: `	-key		"sets struct field [key] (default: defVal)"`,
		},
		"1 fields": {
			flag: &flagData{
				key: "key",
			},
			expected: `	-key		"sets struct field [key]"`,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			gotStr := test.flag.String()

			assert.Equal(t, test.expected, gotStr)
		})
	}
}

func TestGetFlagData(t *testing.T) {
	tests := map[string]struct {
		input    interface{}
		expected *flagData
	}{
		"key": {
			input: struct {
				Name string `flag:"name"`
			}{},
			expected: &flagData{
				key: "name",
			},
		},
		"key & default": {
			input: struct {
				Name string `flag:"name|defVal"`
			}{},
			expected: &flagData{
				key:        "name",
				defaultVal: "defVal",
				usage:      "",
			},
		},
		"key & usage": {
			input: struct {
				Name string `flag:"name||some usage"`
			}{},
			expected: &flagData{
				key:   "name",
				usage: "some usage",
			},
		},
		"key & default & usage": {
			input: struct {
				Name string `flag:"name|defVal|some usage"`
			}{},
			expected: &flagData{
				key:        "name",
				defaultVal: "defVal",
				usage:      "some usage",
			},
		},
		"wrong format": {
			input: struct {
				Name string `flag:"||||"`
			}{},
			expected: nil,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			field := reflect.TypeOf(test.input).Field(0)
			gotFlagData := getFlagData(field)

			assert.Equal(t, test.expected, gotFlagData)
		})
	}
}
