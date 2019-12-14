package configuration

import (
	"reflect"
	"testing"
)

func TestGetTags(t *testing.T) {
	type testStruct struct {
		Name string `json:"jsonVal"  default:"defaultVal"  env:"envVal"  flag:"flagVal"`
	}
	field := reflect.TypeOf(&testStruct{}).Elem().Field(0)

	testCases := []struct {
		name           string
		fn             func(f reflect.StructField) string
		expectedResult string
	}{
		{
			name:           "json",
			fn:             getJSONTag,
			expectedResult: "jsonVal",
		},
		{
			name:           "default",
			fn:             getDefaultTag,
			expectedResult: "defaultVal",
		},
		{
			name:           "env",
			fn:             getEnvTag,
			expectedResult: "envVal",
		},
		{
			name:           "flag",
			fn:             getFlagTag,
			expectedResult: "flagVal",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			actualResult := test.fn(field)
			if !reflect.DeepEqual(test.expectedResult, actualResult) {
				t.Fatalf("\nexpected result: [%s] \nbut got: [%s]", test.expectedResult, actualResult)
			}
		})
	}
}
