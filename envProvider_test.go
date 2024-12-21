// nolint:dupl,paralleltest
package configuration

import (
	"reflect"
	"testing"
)

func TestEnvProvider(t *testing.T) {
	type testStruct struct {
		Name string `env:"ENV_KEY1"`
	}
	testObj := testStruct{}

	fieldType := reflect.TypeOf(&testObj).Elem().Field(0)
	fieldVal := reflect.ValueOf(&testObj).Elem().Field(0)

	provider := NewEnvProvider()
	testValue := "ENV_VAL2"

	t.Setenv("ENV_KEY1", testValue)

	if err := provider.Provide(fieldType, fieldVal); err != nil {
		t.Fatalf("cannot set value: %v", err)
	}

	if !reflect.DeepEqual(testValue, testObj.Name) {
		t.Fatalf("\nexpected result: [%s] \nbut got: [%s]", testValue, testObj.Name)
	}
}

func TestEnvProviderFailed(t *testing.T) {
	type testStruct struct {
		Name string `env:"ENV_KEY1"`
	}
	testObj := testStruct{}

	fieldType := reflect.TypeOf(&testObj).Elem().Field(0)
	fieldVal := reflect.ValueOf(&testObj).Elem().Field(0)

	provider := NewEnvProvider()

	if err := provider.Provide(fieldType, fieldVal); err == nil {
		t.Fatal("must NOT be nil")
	}
}

func TestEnvProvider_EmptyKey(t *testing.T) {
	type testStruct struct {
		Name string `env:""`
	}
	testObj := testStruct{}

	fieldType := reflect.TypeOf(&testObj).Elem().Field(0)
	fieldVal := reflect.ValueOf(&testObj).Elem().Field(0)

	provider := NewEnvProvider()

	err := provider.Provide(fieldType, fieldVal)
	assert(t, "EnvProvider: key is empty", err.Error())
}
