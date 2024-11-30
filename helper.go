package configuration

import (
	"fmt"
	"reflect"
	"strings"
)

func fetchTagKey(t reflect.StructTag) map[string]struct{} {
	keys := map[string]struct{}{}

	pairs := strings.Split(fmt.Sprintf("%s", t), " ")
	if len(pairs) < 1 {
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
