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
	if !reflect.DeepEqual(fieldVal.String(), testStr) {
		t.Fatalf("\nexpected result: [%s] \nbut got: [%s]", testValue, testStr)
	}
}

func TestSetValue_Int8(t *testing.T) {
	var testInt8 int8
	fieldType := reflect.TypeOf(&testInt8).Elem()
	fieldVal := reflect.ValueOf(&testInt8).Elem()
	testValue := "42"

	setValue(fieldType, fieldVal, testValue)
	if fieldVal.Int() != int64(testInt8) {
		t.Fatalf("\nexpected result: [%s] \nbut got: [%d]", testValue, testInt8)
	}
}

func TestSetValue_Uint16(t *testing.T) {
	var testUint16 uint16
	fieldType := reflect.TypeOf(&testUint16).Elem()
	fieldVal := reflect.ValueOf(&testUint16).Elem()
	testValue := "42"

	setValue(fieldType, fieldVal, testValue)
	if fieldVal.Uint() != uint64(testUint16) {
		t.Fatalf("\nexpected result: [%s] \nbut got: [%d]", testValue, testUint16)
	}
}

func TestSetValue_Float32(t *testing.T) {
	var testFloat32 float32
	fieldType := reflect.TypeOf(&testFloat32).Elem()
	fieldVal := reflect.ValueOf(&testFloat32).Elem()
	testValue := "42"

	setValue(fieldType, fieldVal, testValue)
	if fieldVal.Float() != float64(testFloat32) {
		t.Fatalf("\nexpected result: [%s] \nbut got: [%f]", testValue, testFloat32)
	}
}

func TestSetValue_Bool(t *testing.T) {
	var testBool bool
	fieldType := reflect.TypeOf(&testBool).Elem()
	fieldVal := reflect.ValueOf(&testBool).Elem()
	testValue := "true"

	setValue(fieldType, fieldVal, testValue)
	if fieldVal.Bool() != true {
		t.Fatalf("\nexpected result: [%s] \nbut got: [%v]", testValue, testBool)
	}
}

func TestSetPtrValue_Ints(t *testing.T) {
	testValue := "42"

	{
		// int
		var testInt *int
		fieldType := reflect.TypeOf(&testInt).Elem().Elem()
		fieldVal := reflect.ValueOf(&testInt).Elem()

		setPtrValue(fieldType, fieldVal, testValue)
		if testValue != strconv.FormatInt(int64(*testInt), 10) {
			t.Errorf("\nexpected result: [%s] \nbut got: [%v]", testValue, testInt)
		}
	}

	{
		// int8
		var testInt8 *int8
		fieldType := reflect.TypeOf(&testInt8).Elem().Elem()
		fieldVal := reflect.ValueOf(&testInt8).Elem()

		setPtrValue(fieldType, fieldVal, testValue)
		if testValue != strconv.FormatInt(int64(*testInt8), 10) {
			t.Errorf("\nexpected result: [%s] \nbut got: [%v]", testValue, testInt8)
		}
	}

	{
		// int16
		var testInt16 *int16
		fieldType := reflect.TypeOf(&testInt16).Elem().Elem()
		fieldVal := reflect.ValueOf(&testInt16).Elem()

		setPtrValue(fieldType, fieldVal, testValue)
		if testValue != strconv.FormatInt(int64(*testInt16), 10) {
			t.Errorf("\nexpected result: [%s] \nbut got: [%v]", testValue, testInt16)
		}
	}

	{
		// int32
		var testInt32 *int32
		fieldType := reflect.TypeOf(&testInt32).Elem().Elem()
		fieldVal := reflect.ValueOf(&testInt32).Elem()

		setPtrValue(fieldType, fieldVal, testValue)
		if testValue != strconv.FormatInt(int64(*testInt32), 10) {
			t.Errorf("\nexpected result: [%s] \nbut got: [%v]", testValue, testInt32)
		}
	}

	{
		// int64
		var testInt64 *int64
		fieldType := reflect.TypeOf(&testInt64).Elem().Elem()
		fieldVal := reflect.ValueOf(&testInt64).Elem()

		setPtrValue(fieldType, fieldVal, testValue)
		if testValue != strconv.FormatInt(*testInt64, 10) {
			t.Errorf("\nexpected result: [%s] \nbut got: [%v]", testValue, testInt64)
		}
	}
}

func TestSetPtrValue_Uints(t *testing.T) {
	testValue := "42"

	{
		// uint
		var testUint *uint
		fieldType := reflect.TypeOf(&testUint).Elem().Elem()
		fieldVal := reflect.ValueOf(&testUint).Elem()

		setPtrValue(fieldType, fieldVal, testValue)
		if testValue != strconv.FormatUint(uint64(*testUint), 10) {
			t.Errorf("\nexpected result: [%s] \nbut got: [%v]", testValue, testUint)
		}
	}

	{
		// uint8
		var testUint8 *uint8
		fieldType := reflect.TypeOf(&testUint8).Elem().Elem()
		fieldVal := reflect.ValueOf(&testUint8).Elem()

		setPtrValue(fieldType, fieldVal, testValue)
		if testValue != strconv.FormatUint(uint64(*testUint8), 10) {
			t.Errorf("\nexpected result: [%s] \nbut got: [%v]", testValue, testUint8)
		}
	}

	{
		// uint16
		var testUint16 *uint16
		fieldType := reflect.TypeOf(&testUint16).Elem().Elem()
		fieldVal := reflect.ValueOf(&testUint16).Elem()

		setPtrValue(fieldType, fieldVal, testValue)
		if testValue != strconv.FormatUint(uint64(*testUint16), 10) {
			t.Errorf("\nexpected result: [%s] \nbut got: [%v]", testValue, testUint16)
		}
	}

	{
		// uint32
		var testUint32 *uint32
		fieldType := reflect.TypeOf(&testUint32).Elem().Elem()
		fieldVal := reflect.ValueOf(&testUint32).Elem()

		setPtrValue(fieldType, fieldVal, testValue)
		if testValue != strconv.FormatUint(uint64(*testUint32), 10) {
			t.Errorf("\nexpected result: [%s] \nbut got: [%v]", testValue, testUint32)
		}
	}

	{
		// uint64
		var testUint64 *uint64
		fieldType := reflect.TypeOf(&testUint64).Elem().Elem()
		fieldVal := reflect.ValueOf(&testUint64).Elem()

		setPtrValue(fieldType, fieldVal, testValue)
		if testValue != strconv.FormatUint(*testUint64, 10) {
			t.Errorf("\nexpected result: [%s] \nbut got: [%v]", testValue, testUint64)
		}
	}
}

func TestSetPtrValue_Floats(t *testing.T) {
	testValue := "42.0"

	{
		// float32
		var testFloat32 *float32
		fieldType := reflect.TypeOf(&testFloat32).Elem().Elem()
		fieldVal := reflect.ValueOf(&testFloat32).Elem()

		setPtrValue(fieldType, fieldVal, testValue)

		gotStr := strconv.FormatFloat(float64(*testFloat32), 'f', 1, 64)
		if testValue != gotStr {
			t.Errorf("\nexpected result: [%s] \nbut got: [%v]", testValue, gotStr)
		}
	}

	{
		// float64
		var testFloat32 *float64
		fieldType := reflect.TypeOf(&testFloat32).Elem().Elem()
		fieldVal := reflect.ValueOf(&testFloat32).Elem()

		setPtrValue(fieldType, fieldVal, testValue)

		gotStr := strconv.FormatFloat(*testFloat32, 'f', 1, 64)
		if testValue != gotStr {
			t.Errorf("\nexpected result: [%s] \nbut got: [%v]", testValue, gotStr)
		}
	}
}

func TestSetPtrValue_Bool(t *testing.T) {
	var testBool *bool
	fieldType := reflect.TypeOf(&testBool).Elem().Elem()
	fieldVal := reflect.ValueOf(&testBool).Elem()
	testValue := "true"

	setPtrValue(fieldType, fieldVal, testValue)
	if fieldVal.Elem().Bool() != true {
		t.Fatalf("\nexpected result: [%s] \nbut got: [%v]", testValue, testBool)
	}
}

func TestSetValue_StringSlice(t *testing.T) {
	var testStr []string
	fieldType := reflect.TypeOf(&testStr).Elem()
	fieldVal := reflect.ValueOf(&testStr).Elem()
	testValue := "test_val1;test_val2"
	expected := []string{"test_val1", "test_val2"}

	setValue(fieldType, fieldVal, testValue)
	if !reflect.DeepEqual(expected, fieldVal.Interface()) {
		t.Fatalf("\nexpected result: %+v \nbut got: %+v", expected, fieldVal.Interface())
	}
}

func TestSetValue_IntSlice(t *testing.T) {
	var testStr []int
	fieldType := reflect.TypeOf(&testStr).Elem()
	fieldVal := reflect.ValueOf(&testStr).Elem()
	testValue := "1    ; 2 "
	expected := []int{1, 2}

	setValue(fieldType, fieldVal, testValue)
	if !reflect.DeepEqual(expected, fieldVal.Interface()) {
		t.Fatalf("\nexpected result: %+v \nbut got: %+v", expected, fieldVal.Interface())
	}
}

func TestSetValue_FloatSlice(t *testing.T) {
	var testStr []float64
	fieldType := reflect.TypeOf(&testStr).Elem()
	fieldVal := reflect.ValueOf(&testStr).Elem()
	testValue := "1;2.0"
	expected := []float64{1, 2}

	setValue(fieldType, fieldVal, testValue)
	if !reflect.DeepEqual(expected, fieldVal.Interface()) {
		t.Fatalf("\nexpected result: %+v \nbut got: %+v", expected, fieldVal.Interface())
	}
}

func TestSetValue_BoolSlice(t *testing.T) {
	var testStr []bool
	fieldType := reflect.TypeOf(&testStr).Elem()
	fieldVal := reflect.ValueOf(&testStr).Elem()
	testValue := "true; false; "
	expected := []bool{true, false}

	setValue(fieldType, fieldVal, testValue)
	if !reflect.DeepEqual(expected, fieldVal.Interface()) {
		t.Fatalf("\nexpected result: %+v \nbut got: %+v", expected, fieldVal.Interface())
	}
}
