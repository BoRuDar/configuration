package configuration

import (
	"os"
	"reflect"
	"testing"
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

	if !reflect.DeepEqual(testValue, testObj.Name) {
		t.Fatalf("\nexpected result: [%s] \nbut got: [%s]", testValue, testObj.Name)
	}
}
