package configuration

import (
	"reflect"
)

func fetchTagKey(t reflect.StructTag, registered map[string]struct{}) map[string]struct{} {
	keys := map[string]struct{}{}

	for rt := range registered {
		if _, ok := t.Lookup(rt); ok {
			keys[rt] = struct{}{}
		}
	}

	return keys
}
