package configuration

import (
	"flag"
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

func TestFlagProvider_WithErrorHandler(t *testing.T) {
	type testStruct struct {
		Name string `flag:"flag_name4||Description"`
	}
	testObj := testStruct{}
	testValue := "flag_value"
	os.Args = []string{"smth", "-flag_name4=flag_value"}

	fieldType := reflect.TypeOf(&testObj).Elem().Field(0)
	fieldVal := reflect.ValueOf(&testObj).Elem().Field(0)

	eh := func(err error) {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	provider := NewFlagProvider(WithErrorHandler(eh))
	if err := provider.Init(&testObj); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := provider.Provide(fieldType, fieldVal); err != nil {
		t.Fatalf("cannot set value: %v", err)
	}

	assert.Equal(t, testValue, testObj.Name)
}

func TestFlagProvider_WithErrorHandlerAndErr(t *testing.T) {
	type testStruct struct {
		Name string `flag:"flag_name5||||"`
	}
	testObj := testStruct{}
	numberOfExpectedCalls := 1
	callsCounter := 0
	os.Args = []string{""}

	eh := func(err error) {
		callsCounter++

		if err != nil && err.Error() != "flagProvider: wrong flag definition [flag_name5||||]" {
			t.Fatalf("unexpected error")
		}
	}

	if err := NewFlagProvider(WithErrorHandler(eh)).Init(&testObj); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if callsCounter != numberOfExpectedCalls {
		t.Fatalf("error must be called [%d] times but called [%d] times", numberOfExpectedCalls, callsCounter)
	}
}

//func TestFlagProvider_Error(t *testing.T) {
//	type testStruct struct {
//		Name string `flag:"flag_name5||||"`
//	}
//	testObj := testStruct{}
//	os.Args = []string{""}
//
//	eh := func(err error) {
//		if err != nil && err.Error() != ErrNotAPointer.Error() {
//			t.Fatalf("unexpected error: %v", err)
//		}
//	}
//
//	if err := NewFlagProvider(WithErrorHandler(eh)).Init(&testObj); err != nil {
//		t.Fatalf("unexpected error: %v", err)
//	}
//}
