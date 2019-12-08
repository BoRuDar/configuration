package configuration

import (
	"reflect"
	"strconv"
	"testing"
)

func TestSetValue_String(t *testing.T) {
	var testStr string
	fieldType := reflect.TypeOf(&testStr).Elem()
	fieldVal := reflect.ValueOf(&testStr).Elem()
	testValue := "test_val1"

	setValue(fieldType, fieldVal, testValue)
	if !reflect.DeepEqual(testValue, testStr) {
		t.Fatalf("\nexpected result: [%s] \nbut got: [%s]", testValue, testStr)
	}
}

func TestSetValue_Int8(t *testing.T) {
	var testInt8 int8
	fieldType := reflect.TypeOf(&testInt8).Elem()
	fieldVal := reflect.ValueOf(&testInt8).Elem()
	testValue := "42"

	setValue(fieldType, fieldVal, testValue)
	if testValue != strconv.FormatInt(int64(testInt8), 10) {
		t.Fatalf("\nexpected result: [%s] \nbut got: [%d]", testValue, testInt8)
	}
}

func TestSetValue_Uint16(t *testing.T) {
	var testUint16 uint16
	fieldType := reflect.TypeOf(&testUint16).Elem()
	fieldVal := reflect.ValueOf(&testUint16).Elem()
	testValue := "42"

	setValue(fieldType, fieldVal, testValue)
	if testValue != strconv.FormatInt(int64(testUint16), 10) {
		t.Fatalf("\nexpected result: [%s] \nbut got: [%d]", testValue, testUint16)
	}
}

func TestSetValue_Float32(t *testing.T) {
	var testFloat32 float32
	fieldType := reflect.TypeOf(&testFloat32).Elem()
	fieldVal := reflect.ValueOf(&testFloat32).Elem()
	testValue := "42"

	setValue(fieldType, fieldVal, testValue)
	if testValue != strconv.FormatFloat(float64(testFloat32), 'g', -1, 32) {
		t.Fatalf("\nexpected result: [%s] \nbut got: [%f]", testValue, testFloat32)
	}
}

func TestSetValue_Bool(t *testing.T) {
	var testBool bool
	fieldType := reflect.TypeOf(&testBool).Elem()
	fieldVal := reflect.ValueOf(&testBool).Elem()
	testValue := "true"

	setValue(fieldType, fieldVal, testValue)
	if testValue != strconv.FormatBool(true) {
		t.Fatalf("\nexpected result: [%s] \nbut got: [%v]", testValue, testBool)
	}
}
