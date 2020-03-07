package configuration

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"reflect"
	"testing"
)

type testStruct struct {
	Name   string
	Inside struct {
		Beta int
	}
}

func TestFileProvider_yml(t *testing.T) {
	testObj := testStruct{}
	expected := testStruct{
		Inside: struct {
			Beta int
		}{
			Beta: 42,
		},
	}

	provider := NewFileProvider("./testdata/input.yml")

	fieldType := reflect.TypeOf(&testObj).Elem().Field(1).Type.Field(0)
	fieldVal := reflect.ValueOf(&testObj).Elem().Field(1).Field(0)
	fieldPath := []string{"Inside", "Beta"}

	if !provider.Provide(fieldType, fieldVal, fieldPath...) {
		t.Fatal("cannot set value")
	}

	if !reflect.DeepEqual(expected, testObj) {
		t.Fatalf("\nexpected result: [%+v] \nbut got: [%+v]", expected, testObj)
	}
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
		name  string
		input string
	}{
		{
			name:  "json",
			input: "some_name.json",
		},
		{
			name:  "yaml",
			input: "some_name.yaml",
		},
		{
			name:  "yml",
			input: "some_name.yml",
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			gotFn := decodeFunc(test.input)
			if gotFn == nil {
				t.Fatal("expected function but got nil")
			}
		})
	}
}
