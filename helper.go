package configuration

import (
	"reflect"
	"strings"
)

func fetchTagKey(t reflect.StructTag) map[string]struct{} {
	keys := map[string]struct{}{}

	pairs := strings.Split(string(t), " ")
	if len(pairs) == 1 && pairs[0] == "" {
		return keys
	}

	for _, pair := range pairs {
		kv := strings.Split(pair, `:"`)
		if len(kv) < 1 {
			return keys
		}
		keys[kv[0]] = struct{}{}
	}

	return keys
}
