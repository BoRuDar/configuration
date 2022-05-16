package configuration

import (
	"reflect"
	"strings"
	"testing"
)

const fmtMsg = "\nexpected:  (%T)(%v)\ngot:       (%T)(%v)"

func assert(t *testing.T, expected, got any, msg ...string) {
	t.Helper()

	errMsg := fmtMsg
	if len(msg) > 0 {
		errMsg = strings.Join(append(msg, errMsg), " ")
	}

	if !reflect.DeepEqual(expected, got) {
		t.Fatalf(errMsg, expected, expected, got, got)
	}
}
