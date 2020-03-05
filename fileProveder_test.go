package configuration

import (
	"reflect"
	"testing"
)

func TestFileProvider(t *testing.T) {
	type testStruct struct {
		Name   string
		Inside struct {
			Beta int
		}
	}

	expected := testStruct{
		Name: "test_name_json",
		Inside: struct {
			Beta int
		}{
			Beta: 42,
		},
	}

	testObj := testStruct{}
	provider := NewFileProvider("./testdata/input.yml")

	fieldType := reflect.TypeOf(&testObj).Elem().Field(0)
	fieldVal := reflect.ValueOf(&testObj).Elem().Field(0)

	if !provider.Provide(fieldType, fieldVal) {
		t.Fatal("cannot set value")
	}

	if !reflect.DeepEqual(expected, testObj) {
		t.Fatalf("\nexpected result: [%+v] \nbut got: [%+v]", expected, testObj)
	}
}
