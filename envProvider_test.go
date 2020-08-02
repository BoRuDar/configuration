package configuration

import (
	"os"
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

	removeEnvKey, err := setEnv("ENV_KEY1", testValue)
	if err != nil {
		t.Fatal("unexpected err: ", err)
	}
	defer removeEnvKey()

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

func setEnv(key, val string) (func(), error) {
	return func() {
		_ = os.Unsetenv(key)
	}, os.Setenv(key, val)
}
