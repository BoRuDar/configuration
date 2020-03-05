package configuration

import (
	"reflect"
	"testing"
)

func TestFileProvider_yml(t *testing.T) {
	type testStruct struct {
		Inside struct {
			Beta int
		}
	}

	expected := testStruct{
		Inside: struct {
			Beta int
		}{
			Beta: 42,
		},
	}

	testObj := testStruct{}
	provider := NewFileProvider("./testdata/input.yml")
	fieldPath := []string{"Inside", "Beta"}

	fieldType := reflect.TypeOf(&testObj).Elem().Field(0).Type.Field(0)
	fieldVal := reflect.ValueOf(&testObj).Elem().Field(0).Field(0)

	if !provider.Provide(fieldType, fieldVal, fieldPath...) {
		t.Fatal("cannot set value")
	}

	if !reflect.DeepEqual(expected, testObj) {
		t.Fatalf("\nexpected result: [%+v] \nbut got: [%+v]", expected, testObj)
	}
}
