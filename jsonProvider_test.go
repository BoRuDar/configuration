// nolint:paralleltest
package configuration

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestJSONFileProvider_json(t *testing.T) {
	type test struct {
		Timeout time.Duration `file_json:"timeout"`
	}

	testObj := test{}
	expected := test{
		Timeout: time.Millisecond * 101,
	}

	fieldType := reflect.TypeOf(&testObj).Elem().Field(0)
	fieldVal := reflect.ValueOf(&testObj).Elem().Field(0)

	p := NewJSONFileProvider("./testdata/input.json")
	if err := p.Init(&testObj); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := p.Provide(fieldType, fieldVal); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert(t, expected, testObj)
}

// nolint:errchkjson,musttag
func TestFindValStrByPath(t *testing.T) {
	type embedded struct {
		Beta int `file_json:"inside.beta"`
	}

	type testStruct struct {
		Name    string        `file_json:"name"`
		Timeout time.Duration `file_json:"timeout"`
		Inside  embedded
	}

	var testObjFromJSON any
	data, _ := json.Marshal(testStruct{
		Name:   "test",
		Inside: embedded{Beta: 42},
	})
	_ = json.Unmarshal(data, &testObjFromJSON)

	tests := []struct {
		name         string
		input        any
		path         []string
		expectedStr  string
		expectedBool bool
	}{
		{
			name:         "empty path",
			input:        nil,
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
			name:         "substructures | Inside.Beta | json",
			input:        testObjFromJSON,
			path:         []string{"Inside", "Beta"},
			expectedStr:  "42",
			expectedBool: true,
		},
		{
			name:         "not found",
			input:        testObjFromJSON,
			path:         []string{"notfound"},
			expectedStr:  "",
			expectedBool: false,
		},
	}

	for _, tt := range tests {
		test := tt

		t.Run(test.name, func(t *testing.T) {
			gotStr, gotBool := findValStrByPath(test.input, test.path)
			if gotStr != test.expectedStr || gotBool != test.expectedBool {
				t.Fatalf("expected: [%q %v] but got [%q %v]", test.expectedStr, test.expectedBool, gotStr, gotBool)
			}
		})
	}
}

func TestFileProvider_Init(t *testing.T) {
	type Cfg struct {
		Test int `file_json:"void."`
	}

	_, err := New[Cfg](NewJSONFileProvider("./testdata/dummy.file"))
	assert(t, "cannot init [JSONFileProvider] provider: file must have .json extension", err.Error())

	_, err = New[Cfg](NewJSONFileProvider("./testdata/malformed_input.json"))
	assert(t, "cannot init [JSONFileProvider] provider: JSONFileProvider.Init: invalid character '}' looking for beginning of value", err.Error())
}

func TestJSONProvider_empty_tag(t *testing.T) {
	type testStruct struct {
		Test int `file_json:""`
	}

	testObj := testStruct{}

	fieldType := reflect.TypeOf(&testObj).Elem().Field(0)
	fieldVal := reflect.ValueOf(&testObj).Elem().Field(0)

	provider := NewJSONFileProvider("./testdata/input.json")
	err := provider.Provide(fieldType, fieldVal)
	assert(t, "JSONFileProvider: key is empty", err.Error())
}
