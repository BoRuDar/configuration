package configuration

import (
	"reflect"
	"testing"
)

func TestDefaultProvider(t *testing.T) {
	type testStruct struct {
		Name string `default:"default_provider_val"`
	}
	testObj := testStruct{}

	fieldType := reflect.TypeOf(&testObj).Elem().Field(0)
	fieldVal := reflect.ValueOf(&testObj).Elem().Field(0)

	provider := NewDefaultProvider()
	testValue := "default_provider_val"

	if !provider.Provide(fieldType, fieldVal) {
		t.Fatal("cannot set value")
	}

	if !reflect.DeepEqual(testValue, testObj.Name) {
		t.Fatalf("\nexpected result: [%s] \nbut got: [%s]", testValue, testObj.Name)
	}
}

func TestDefaultProviderFailed(t *testing.T) {
	type testStruct struct {
		Name string
	}
	testObj := testStruct{}

	fieldType := reflect.TypeOf(&testObj).Elem().Field(0)
	fieldVal := reflect.ValueOf(&testObj).Elem().Field(0)

	provider := NewDefaultProvider()

	if provider.Provide(fieldType, fieldVal) {
		t.Fatal("must be false")
	}
}
