// nolint:dupl
package configuration

import (
	"reflect"
	"testing"
)

func TestDefaultProvider(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Name string `default:"default_provider_val"`
	}
	testObj := testStruct{}

	fieldType := reflect.TypeOf(&testObj).Elem().Field(0)
	fieldVal := reflect.ValueOf(&testObj).Elem().Field(0)

	provider := NewDefaultProvider()
	testValue := "default_provider_val"

	if err := provider.Provide(fieldType, fieldVal); err != nil {
		t.Fatal("cannot set value")
	}

	if !reflect.DeepEqual(testValue, testObj.Name) {
		t.Fatalf("\nexpected result: [%s] \nbut got: [%s]", testValue, testObj.Name)
	}
}

func TestDefaultProviderPtr(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Name *string `default:"default_provider_val_ptr"`
	}
	testObj := testStruct{}

	fieldType := reflect.TypeOf(&testObj).Elem().Field(0)
	fieldVal := reflect.ValueOf(&testObj).Elem().Field(0)

	provider := NewDefaultProvider()
	testValue := "default_provider_val_ptr"

	if err := provider.Provide(fieldType, fieldVal); err != nil {
		t.Fatal("cannot set value")
	}

	if !reflect.DeepEqual(testValue, *testObj.Name) {
		t.Fatalf("\nexpected result: [%s] \nbut got: [%s]", testValue, *testObj.Name)
	}
}

func TestDefaultProviderFailed(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Name string `default:""`
	}
	testObj := testStruct{}

	fieldType := reflect.TypeOf(&testObj).Elem().Field(0)
	fieldVal := reflect.ValueOf(&testObj).Elem().Field(0)

	provider := NewDefaultProvider()

	if err := provider.Provide(fieldType, fieldVal); err == nil {
		t.Fatal("must not be nil")
	}
}

func TestDefaultProviderFailedNoTag(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Name string
	}
	testObj := testStruct{}

	fieldType := reflect.TypeOf(&testObj).Elem().Field(0)
	fieldVal := reflect.ValueOf(&testObj).Elem().Field(0)

	provider := NewDefaultProvider()

	if err := provider.Provide(fieldType, fieldVal); err != nil {
		t.Fatal("must be nil")
	}
}
