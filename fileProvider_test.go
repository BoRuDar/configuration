package configuration

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

type testStruct struct {
	Name   string
	Inside struct {
		Beta int
	}
	Timeout time.Duration
}

func TestFileProvider_yml(t *testing.T) {
	testObj := testStruct{}
	expected := testStruct{
		Inside: struct {
			Beta int
		}{
			Beta: 42,
		},
		Timeout: time.Millisecond * 100,
	}

	provider, err := NewFileProvider("./testdata/input.yml")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	var ( // field: Inside.Beta
		fieldType = reflect.TypeOf(&testObj).Elem().Field(1).Type.Field(0)
		fieldVal  = reflect.ValueOf(&testObj).Elem().Field(1).Field(0)
		fieldPath = []string{"Inside", "Beta"}
	)
	var ( // field: Timeout
		fieldType2 = reflect.TypeOf(&testObj).Elem().Field(2)
		fieldVal2  = reflect.ValueOf(&testObj).Elem().Field(2)
		fieldPath2 = []string{"Timeout"}
	)

	err2 := provider.Provide(fieldType, fieldVal, fieldPath...)
	err3 := provider.Provide(fieldType2, fieldVal2, fieldPath2...)

	assert.Nil(t, err2, "cannot set value for Inside.Beta")
	assert.Nil(t, err3, "cannot set value for Timeout")
	assert.Equal(t, expected, testObj)
}

func TestFileProvider_json(t *testing.T) {
	testObj := testStruct{}
	expected := testStruct{
		Timeout: time.Millisecond * 101,
	}

	provider, err := NewFileProvider("./testdata/input.json")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	fieldType := reflect.TypeOf(&testObj).Elem().Field(2)
	fieldVal := reflect.ValueOf(&testObj).Elem().Field(2)
	fieldPath := []string{"Timeout"}

	err1 := provider.Provide(fieldType, fieldVal, fieldPath...)

	assert.Nil(t, err1, "cannot set value")
	assert.Equal(t, expected, testObj)
}

func TestFindValStrByPath(t *testing.T) {
	var testObjFromYAML interface{}
	data, _ := yaml.Marshal(testStruct{
		Name:   "test",
		Inside: struct{ Beta int }{Beta: 42},
	})
	_ = yaml.Unmarshal(data, &testObjFromYAML)

	var testObjFromJSON interface{}
	data, _ = json.Marshal(testStruct{
		Name:   "test",
		Inside: struct{ Beta int }{Beta: 42},
	})
	_ = json.Unmarshal(data, &testObjFromJSON)

	tests := []struct {
		name         string
		input        interface{}
		path         []string
		expectedStr  string
		expectedBool bool
	}{
		{
			name:         "empty path",
			path:         nil,
			expectedStr:  "",
			expectedBool: false,
		},
		{
			name:         "at root level | Name | json",
			input:        testObjFromJSON,
			path:         []string{"Name"},
			expectedStr:  "test",
			expectedBool: true,
		},
		{
			name:         "at root level | Name | yaml",
			input:        testObjFromYAML,
			path:         []string{"Name"},
			expectedStr:  "test",
			expectedBool: true,
		},
		{
			name:         "substructures | Inside.Beta | json",
			input:        testObjFromJSON,
			path:         []string{"Inside", "Beta"},
			expectedStr:  "42",
			expectedBool: true,
		},
		{
			name:         "substructures | Inside.Beta | yaml",
			input:        testObjFromYAML,
			path:         []string{"Inside", "Beta"},
			expectedStr:  "42",
			expectedBool: true,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			gotStr, gotBool := findValStrByPath(testObjFromYAML, test.path)
			if gotStr != test.expectedStr || gotBool != test.expectedBool {
				t.Fatalf("expected: [%q %v] but got [%q %v]", test.expectedStr, test.expectedBool, gotStr, gotBool)
			}
		})
	}
}

func TestDecodeFunc(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "json",
			input:       "some_name.json",
			expectError: false,
		},
		{
			name:        "yaml",
			input:       "some_name.yaml",
			expectError: false,
		},
		{
			name:        "yml",
			input:       "some_name.yml",
			expectError: false,
		},
		{
			name:        "unsupported",
			input:       "some_name.bad",
			expectError: true,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			_, err := decodeFunc(test.input)
			if err != nil && !test.expectError {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
