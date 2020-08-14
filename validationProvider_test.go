package configuration

import (
	"reflect"
	"testing"
)

func TestValidationProvider(t *testing.T) {
	type testStruct struct {
		Name string `validate:"required"`
	}
	testObj := testStruct{}

	fieldType := reflect.TypeOf(&testObj).Elem().Field(0)
	fieldVal := reflect.ValueOf(&testObj).Elem().Field(0)

	provider := NewValidationProvider(NewDefaultProvider())
	testValue := "required"

	if err := provider.Provide(fieldType, fieldVal); err != nil {
		t.Fatal("cannot set value")
	}

	if !reflect.DeepEqual(testValue, testObj.Name) {
		t.Fatalf("\nexpected result: [%s] \nbut got: [%s]", testValue, testObj.Name)
	}
}