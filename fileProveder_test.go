package configuration

import (
	"os"
	"reflect"
	"testing"
)

func TestFileProvider(t *testing.T) {
	type testStruct struct {
		Name string `json:"key_name"`
	}
	expectedValue := "test_val"

	file, err := os.Open("./json_sample.json")
	if err != nil {
		t.Fatalf("cannot open test file: %v", err)
	}
	defer file.Close()

	testObj := testStruct{}
	provider := NewFileProvider(&testObj, file)

	fieldType := reflect.TypeOf(&testObj).Elem().Field(0)
	fieldVal := reflect.ValueOf(&testObj).Elem().Field(0)

	if !provider.Provide(fieldType, fieldVal) {
		t.Fatal("cannot set value")
	}

	if !reflect.DeepEqual(expectedValue, testObj.Name) {
		t.Fatalf("\nexpected result: [%s] \nbut got: [%s]", expectedValue, testObj.Name)
	}
}
