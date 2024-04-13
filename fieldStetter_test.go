// nolint:dupl,goconst
package configuration

import (
	"net"
	"reflect"
	"strings"
	"testing"
	"time"
)

// SetValue tests

func TestSetValue_String(t *testing.T) {
	t.Parallel()

	var testStr string
	fieldType := reflect.TypeOf(&testStr).Elem()
	fieldVal := reflect.ValueOf(&testStr).Elem()
	testValue := "test_val1"
	expectedValue := "test_val1"

	err := setValue(fieldType, fieldVal, testValue)
	assert(t, nil, err)
	assert(t, expectedValue, testStr)
}

func TestSetValue_Int8(t *testing.T) {
	t.Parallel()

	var testInt8 int8
	fieldType := reflect.TypeOf(&testInt8).Elem()
	fieldVal := reflect.ValueOf(&testInt8).Elem()
	testValue := "42"
	expectedValue := int8(42)

	err := setValue(fieldType, fieldVal, testValue)
	assert(t, nil, err)
	assert(t, expectedValue, testInt8)
}

func TestSetValue_Uint16(t *testing.T) {
	t.Parallel()

	var testUint16 uint16
	fieldType := reflect.TypeOf(&testUint16).Elem()
	fieldVal := reflect.ValueOf(&testUint16).Elem()
	testValue := "42"
	expectedValue := uint16(42)

	err := setValue(fieldType, fieldVal, testValue)
	assert(t, nil, err)
	assert(t, expectedValue, testUint16)
}

func TestSetValue_Int64(t *testing.T) {
	t.Parallel()

	var (
		testInt64 int64
		fieldVal  = reflect.ValueOf(&testInt64).Elem()
		testValue = "42"
	)

	setInt64(fieldVal, testValue)
	assert(t, testInt64, testInt64)
}

func TestSetValue_Duration(t *testing.T) {
	t.Parallel()

	var (
		testDuration time.Duration
		fieldVal     = reflect.ValueOf(&testDuration).Elem()
		testValue    = "42ms"
		expectedVal  = time.Millisecond * 42
	)

	setInt64(fieldVal, testValue)
	assert(t, expectedVal, testDuration)
}

func TestSetValue_Float(t *testing.T) {
	t.Parallel()

	var testFloat32 float32
	fieldType := reflect.TypeOf(&testFloat32).Elem()
	fieldVal := reflect.ValueOf(&testFloat32).Elem()
	testValue := "42"
	expectedValue := float32(42.0)

	err := setValue(fieldType, fieldVal, testValue)
	assert(t, nil, err)
	assert(t, expectedValue, testFloat32)
}

func TestSetValue_Bool(t *testing.T) {
	t.Parallel()

	var testBool bool
	fieldType := reflect.TypeOf(&testBool).Elem()
	fieldVal := reflect.ValueOf(&testBool).Elem()
	testValue := "true"
	expectedValue := true

	err := setValue(fieldType, fieldVal, testValue)
	assert(t, nil, err)
	assert(t, expectedValue, testBool)
}

// SetPtr tests

func TestSetPtr_String(t *testing.T) {
	t.Parallel()

	var testStr *string
	fieldType := reflect.TypeOf(&testStr).Elem()
	fieldVal := reflect.ValueOf(&testStr).Elem()
	testValue := "test_val1"
	expected := ToPtr[string]("test_val1")

	err := setPtrValue(fieldType, fieldVal, testValue)
	assert(t, nil, err)
	assert(t, expected, testStr)
}

func TestSetPtrValue_Ints(t *testing.T) {
	t.Parallel()

	testValue := "42"

	{
		var testInt *int
		fieldType := reflect.TypeOf(&testInt).Elem()
		fieldVal := reflect.ValueOf(&testInt).Elem()
		expectedVal := ToPtr[int](42)

		err := setPtrValue(fieldType, fieldVal, testValue)
		assert(t, nil, err)
		assert(t, expectedVal, testInt)
	}

	{
		var testInt8 *int8
		fieldType := reflect.TypeOf(&testInt8).Elem()
		fieldVal := reflect.ValueOf(&testInt8).Elem()
		expectedVal := ToPtr[int8](42)

		err := setPtrValue(fieldType, fieldVal, testValue)
		assert(t, nil, err)
		assert(t, expectedVal, testInt8)
	}

	{
		var testInt16 *int16
		fieldType := reflect.TypeOf(&testInt16).Elem()
		fieldVal := reflect.ValueOf(&testInt16).Elem()
		expectedVal := ToPtr[int16](42)

		err := setPtrValue(fieldType, fieldVal, testValue)
		assert(t, nil, err)
		assert(t, expectedVal, testInt16)
	}

	{
		var testInt32 *int32
		fieldType := reflect.TypeOf(&testInt32).Elem()
		fieldVal := reflect.ValueOf(&testInt32).Elem()
		expectedVal := ToPtr[int32](42)

		err := setPtrValue(fieldType, fieldVal, testValue)
		assert(t, nil, err)
		assert(t, expectedVal, testInt32)
	}

	{
		var testInt64 *int64
		fieldType := reflect.TypeOf(&testInt64).Elem()
		fieldVal := reflect.ValueOf(&testInt64).Elem()
		expectedVal := ToPtr[int64](42)

		err := setPtrValue(fieldType, fieldVal, testValue)
		assert(t, nil, err)
		assert(t, expectedVal, testInt64)
	}
}

func TestSetPtrValue_Uints(t *testing.T) {
	t.Parallel()

	testValue := "42"

	{
		var testUint *uint
		fieldType := reflect.TypeOf(&testUint).Elem()
		fieldVal := reflect.ValueOf(&testUint).Elem()
		expectedVal := ToPtr[uint](42)

		err := setPtrValue(fieldType, fieldVal, testValue)
		assert(t, nil, err)
		assert(t, expectedVal, testUint)
	}

	{
		var testUint8 *uint8
		fieldType := reflect.TypeOf(&testUint8).Elem()
		fieldVal := reflect.ValueOf(&testUint8).Elem()
		expectedVal := ToPtr[uint8](42)

		err := setPtrValue(fieldType, fieldVal, testValue)
		assert(t, nil, err)
		assert(t, expectedVal, testUint8)
	}

	{
		var testUint16 *uint16
		fieldType := reflect.TypeOf(&testUint16).Elem()
		fieldVal := reflect.ValueOf(&testUint16).Elem()
		expectedVal := ToPtr[uint16](42)

		err := setPtrValue(fieldType, fieldVal, testValue)
		assert(t, nil, err)
		assert(t, expectedVal, testUint16)
	}

	{
		var testUint32 *uint32
		fieldType := reflect.TypeOf(&testUint32).Elem()
		fieldVal := reflect.ValueOf(&testUint32).Elem()
		expectedVal := ToPtr[uint32](42)

		err := setPtrValue(fieldType, fieldVal, testValue)
		assert(t, nil, err)
		assert(t, expectedVal, testUint32)
	}

	{
		var testUint64 *uint64
		fieldType := reflect.TypeOf(&testUint64).Elem()
		fieldVal := reflect.ValueOf(&testUint64).Elem()
		expectedVal := ToPtr[uint64](42)

		err := setPtrValue(fieldType, fieldVal, testValue)
		assert(t, nil, err)
		assert(t, expectedVal, testUint64)
	}
}

func TestSetPtrValue_Floats(t *testing.T) {
	t.Parallel()

	testValue := "42.0"

	{
		var testFloat32 *float32
		fieldType := reflect.TypeOf(&testFloat32).Elem()
		fieldVal := reflect.ValueOf(&testFloat32).Elem()
		expectedVal := ToPtr[float32](42)

		err := setPtrValue(fieldType, fieldVal, testValue)
		assert(t, nil, err)
		assert(t, expectedVal, testFloat32)
	}

	{
		var testFloat64 *float64
		fieldType := reflect.TypeOf(&testFloat64).Elem()
		fieldVal := reflect.ValueOf(&testFloat64).Elem()
		expectedVal := ToPtr[float64](42)

		err := setPtrValue(fieldType, fieldVal, testValue)
		assert(t, nil, err)
		assert(t, expectedVal, testFloat64)
	}
}

func TestSetPtrValue_Bool(t *testing.T) {
	t.Parallel()

	var testBool *bool
	fieldType := reflect.TypeOf(&testBool).Elem()
	fieldVal := reflect.ValueOf(&testBool).Elem()
	testValue := "true"
	expectedVal := ToPtr[bool](true)

	err := setPtrValue(fieldType, fieldVal, testValue)
	assert(t, nil, err)
	assert(t, expectedVal, testBool)
}

// SetValue slice tests

func TestSetValue_StringSlice(t *testing.T) {
	t.Parallel()

	var testStr []string
	fieldType := reflect.TypeOf(&testStr).Elem()
	fieldVal := reflect.ValueOf(&testStr).Elem()
	testValue := "test_val1;test_val2"
	expected := []string{"test_val1", "test_val2"}

	err := setValue(fieldType, fieldVal, testValue)
	assert(t, nil, err)
	assert(t, expected, fieldVal.Interface())
}

func TestSetValue_StringSliceSingleElement(t *testing.T) {
	t.Parallel()

	var testStr []string
	fieldType := reflect.TypeOf(&testStr).Elem()
	fieldVal := reflect.ValueOf(&testStr).Elem()
	testValue := "test_val1"
	expected := []string{"test_val1"}

	err := setValue(fieldType, fieldVal, testValue)
	assert(t, nil, err)
	assert(t, expected, fieldVal.Interface())
}

func TestSetValue_IntSlice(t *testing.T) {
	t.Parallel()

	var testStr []int
	fieldType := reflect.TypeOf(&testStr).Elem()
	fieldVal := reflect.ValueOf(&testStr).Elem()
	testValue := "1    ; 2 "
	expected := []int{1, 2}

	err := setValue(fieldType, fieldVal, testValue)
	assert(t, nil, err)
	assert(t, expected, fieldVal.Interface())
}

func TestSetValue_UintSlice(t *testing.T) {
	t.Parallel()

	var (
		testStr   []uint
		fieldType = reflect.TypeOf(&testStr).Elem()
		fieldVal  = reflect.ValueOf(&testStr).Elem()
		testValue = "1  ; 2 "
		expected  = []uint{1, 2}
	)

	err := setValue(fieldType, fieldVal, testValue)
	assert(t, nil, err)
	assert(t, expected, fieldVal.Interface())
}

func TestSetValue_FloatSlice(t *testing.T) {
	t.Parallel()

	var testStr []float64
	fieldType := reflect.TypeOf(&testStr).Elem()
	fieldVal := reflect.ValueOf(&testStr).Elem()
	testValue := "1;2.0"
	expected := []float64{1, 2}

	err := setValue(fieldType, fieldVal, testValue)
	assert(t, nil, err)
	assert(t, expected, fieldVal.Interface())
}

func TestSetValue_BoolSlice(t *testing.T) {
	t.Parallel()

	var testStr []bool
	fieldType := reflect.TypeOf(&testStr).Elem()
	fieldVal := reflect.ValueOf(&testStr).Elem()
	testValue := "true; false; "
	expected := []bool{true, false}

	err := setValue(fieldType, fieldVal, testValue)
	assert(t, nil, err)
	assert(t, expected, testStr)
}

func TestSetValue_EmptySlice(t *testing.T) {
	t.Parallel()

	var testStr []bool
	fieldType := reflect.TypeOf(&testStr).Elem()
	fieldVal := reflect.ValueOf(&testStr).Elem()
	testValue := " "

	err := setValue(fieldType, fieldVal, testValue)
	assert(t, "setSlice: got empty slice", err.Error())
}

func TestSetValue_Unsupported(t *testing.T) {
	t.Parallel()

	var testStr chan struct{}
	fieldType := reflect.TypeOf(&testStr).Elem()
	fieldVal := reflect.ValueOf(&testStr).Elem()
	testValue := "true; false; "

	err := setValue(fieldType, fieldVal, testValue)
	assert(t, "setValue: unsupported type: chan", err.Error())

	err = setPtrValue(fieldType, fieldVal, testValue)
	assert(t, "setPtrValue: unsupported type: chan", err.Error())

	err = setSlice(fieldType, fieldVal, testValue)
	assert(t, "setSlice: unsupported type of slice item: struct", err.Error())
}

func TestSetValue_IntPtrSlice(t *testing.T) {
	t.Parallel()

	var testIntSlice []*int
	fieldType := reflect.TypeOf(&testIntSlice).Elem()
	fieldVal := reflect.ValueOf(&testIntSlice).Elem()
	testValue := "1;2;3"
	expected := []*int{ToPtr(1), ToPtr(2), ToPtr(3)}

	err := setValue(fieldType, fieldVal, testValue)
	assert(t, nil, err)
	assert(t, expected, testIntSlice)
}

func TestSetValue_IntPtrSlice_Err(t *testing.T) {
	t.Parallel()

	var testStr []*struct{}
	fieldType := reflect.TypeOf(&testStr).Elem()
	fieldVal := reflect.ValueOf(&testStr).Elem()
	testValue := "1;2;4"

	err := setValue(fieldType, fieldVal, testValue)
	if err == nil {
		t.Fatal("expected err but got nil")
	}

	if err.Error() != "setSlice: cannot set type [*struct {}] at index [0]" {
		t.Fatalf("wrong error: %v", err)
	}
}

type testCfgSetField struct {
	HostOne *ipTest `default:"127.0.0.1"`
	HostTwo ipTest  `default:"127.0.0.2"`
	Hosts   ipsTest `default:"10.0.0.1,10.0.0.2"`
	NameOne string  `default:"one"`
	NameTwo *string `default:"two"`
}

type ipTest net.IP

func (it *ipTest) SetField(_ reflect.StructField, val reflect.Value, valStr string) error {
	i := ipTest(net.ParseIP(valStr))

	if val.Kind() == reflect.Pointer {
		val.Set(reflect.ValueOf(&i))
	} else {
		val.Set(reflect.ValueOf(i))
	}

	return nil
}

type ipsTest []ipTest

func (it *ipsTest) SetField(sf reflect.StructField, val reflect.Value, valStr string) error {
	var (
		strIPs = strings.Split(valStr, ",")
		ips    = make(ipsTest, len(strIPs))
	)

	for i, ip := range strIPs {
		if err := ips[i].SetField(sf, reflect.ValueOf(&ips[i]).Elem(), ip); err != nil {
			return err
		}
	}

	if val.Kind() == reflect.Pointer {
		val.Set(reflect.ValueOf(&ips))
	} else {
		val.Set(reflect.ValueOf(ips))
	}

	return nil
}

func Test_CustomFieldSetter(t *testing.T) {
	t.Parallel()

	var cfg testCfgSetField
	err := FromEnvAndDefault(&cfg)

	assert(t, nil, err)
	assert(t, "127.0.0.1", net.IP(*cfg.HostOne).String())
	assert(t, "127.0.0.2", net.IP(cfg.HostTwo).String())
	assert(t, "one", cfg.NameOne)
	assert(t, "two", *cfg.NameTwo)
	assert(t, ipsTest([]ipTest{
		ipTest(net.ParseIP("10.0.0.1")),
		ipTest(net.ParseIP("10.0.0.2")),
	}), cfg.Hosts)
}
