package configuration

import (
	"flag"
	"fmt"
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
	testValue := "flag_value"

	fieldType := reflect.TypeOf(&testObj).Elem().Field(0)
	fieldVal := reflect.ValueOf(&testObj).Elem().Field(0)

	provider := NewFlagProvider()

	if err := provider.Init(&testObj); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := provider.Provide(fieldType, fieldVal); err != nil {
		t.Fatalf("cannot set value: %v", err)
	}

	assert.Equal(t, testValue, testObj.Name)
}

func TestFlagProvider_WithDescription(t *testing.T) {
	type testStruct struct {
		Name string `flag:"flag_name2||Description"`
	}
	testObj := testStruct{}
	testValue := "flag_value"
	os.Args = []string{"smth", "-flag_name2=flag_value"}

	fieldType := reflect.TypeOf(&testObj).Elem().Field(0)
	fieldVal := reflect.ValueOf(&testObj).Elem().Field(0)

	provider := NewFlagProvider()

	if err := provider.Init(&testObj); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := provider.Provide(fieldType, fieldVal); err != nil {
		t.Fatalf("cannot set value: %v", err)
	}

	assert.Equal(t, testValue, testObj.Name)
}

func TestFlagProvider_WithDefault(t *testing.T) {
	type testStruct struct {
		Name string `flag:"flag_name3|default_val"`
	}
	testObj := testStruct{}
	testValue := "default_val"

	fieldType := reflect.TypeOf(&testObj).Elem().Field(0)
	fieldVal := reflect.ValueOf(&testObj).Elem().Field(0)

	provider := NewFlagProvider()
	if err := provider.Init(&testObj); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := provider.Provide(fieldType, fieldVal); err != nil {
		t.Fatalf("cannot set value: %v", err)
	}

	assert.Equal(t, testValue, testObj.Name)
}

func TestGetFlagData(t *testing.T) {
	tests := map[string]struct {
		input    interface{}
		expected *flagData
		hasErr   bool
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
			hasErr:   true,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			field := reflect.TypeOf(test.input).Field(0)
			gotFlagData, err := getFlagData(field)

			assert.Equal(t, test.hasErr, err != nil)
			assert.Equal(t, test.expected, gotFlagData)
		})
	}
}

func TestFlagProvider_CustomFlagSet(t *testing.T) {
	type testStruct struct {
		Name string `flag:"flag_name3||Description"`
	}
	testObj := testStruct{}
	testValue := "flag_value"
	os.Args = []string{"smth", "-flag_name3=flag_value"}

	fieldType := reflect.TypeOf(&testObj).Elem().Field(0)
	fieldVal := reflect.ValueOf(&testObj).Elem().Field(0)

	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	provider := NewFlagProvider(WithFlagSet(fs))

	if err := provider.Init(&testObj); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := provider.Provide(fieldType, fieldVal); err != nil {
		t.Fatalf("cannot set value: %v", err)
	}

	assert.Equal(t, testValue, testObj.Name)
}

func TestFlagProvider_Panic(t *testing.T) {
	testObj := struct {
		s struct{}
	}{}

	err := NewFlagProvider().Init(&testObj)
	assert.Error(t, err)
	assert.Equal(t, "got panic: reflect.Value.Interface: cannot return value obtained from unexported field or method", err.Error())
}

func TestFlagProvider_ErrNotAPointer(t *testing.T) {
	type testStruct struct {
		Name string `flag:"flag_name6||||"`
	}
	testObj := testStruct{}
	os.Args = []string{""}

	if err := NewFlagProvider().Init(testObj); err != ErrNotAPointer {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFlagProvider_Errors(t *testing.T) {
	testCases := map[string]struct {
		obj           interface{}
		initErr       error
		providerError error
	}{
		"Empty value": {
			obj: &struct {
				Name string `flag:"flag_name7||Description"`
			}{},
			providerError: ErrEmptyValue,
		},
		"Tag is not unique": {
			obj: &struct {
				Name  string `flag:"flag_name8"`
				Name2 string `flag:"flag_name8"`
			}{},
			initErr: fmt.Errorf("%w: flag_name8", ErrTagNotUnique),
		},
		"No tag": {
			obj: &struct {
				Name string
			}{},
			providerError: ErrNoTag,
		},
	}

	for name, test := range testCases {
		test := test

		t.Run(name, func(t *testing.T) {
			fieldType := reflect.TypeOf(test.obj).Elem().Field(0)
			fieldVal := reflect.ValueOf(test.obj).Elem().Field(0)

			provider := NewFlagProvider()
			if err := provider.Init(test.obj); err != nil {
				if test.initErr != nil && err.Error() == test.initErr.Error() {
					return
				}

				t.Fatalf("unexpected init error: %v", err)
			}

			if err := provider.Provide(fieldType, fieldVal); err != nil {
				if test.providerError != nil && err.Error() == test.providerError.Error() {
					return
				}

				t.Fatalf("unexpected provider error: %v", err)
			}
		})
	}
}
