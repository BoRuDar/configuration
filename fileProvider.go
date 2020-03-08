package configuration

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
)

// NewFileProvider creates new provider which read values from files (json, yaml)
func NewFileProvider(fileName string) (fp fileProvider) {
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	if b, err := ioutil.ReadAll(file); err == nil {
		if fn := decodeFunc(fileName); fn != nil {
			err := fn(b, &fp.fileData)
			if err != nil {
				log.Println(err)
			}
		}
	}
	return
}

type fileProvider struct {
	fileData interface{}
}

func (fp fileProvider) Provide(field reflect.StructField, v reflect.Value, path ...string) bool {
	valStr, ok := findValStrByPath(fp.fileData, path)
	if !ok {
		return false
	}

	SetField(field, v, valStr)
	logf("fileProvider: set [%s] to field [%s]", valStr, strings.Join(path, "."))
	return true
}

func decodeFunc(fileName string) func(data []byte, v interface{}) error {
	fileName = strings.ToLower(fileName)

	if strings.HasSuffix(fileName, ".json") {
		return json.Unmarshal
	}
	if strings.HasSuffix(fileName, ".yaml") {
		return yaml.Unmarshal
	}
	if strings.HasSuffix(fileName, ".yml") {
		return yaml.Unmarshal
	}

	logf("fileProvider: unsupported file type: %q", fileName)
	return nil
}

func findValStrByPath(i interface{}, path []string) (string, bool) {
	if len(path) == 0 {
		return "", false
	}
	firstInPath := strings.ToLower(path[0])

	currentFieldStr, ok := i.(map[string]interface{}) // unmarshaled from json
	if !ok {
		currentFieldIface, ok := i.(map[interface{}]interface{}) // unmarshaled from yaml
		if !ok {
			return "", false
		}

		currentFieldStr = map[string]interface{}{}
		for k, v := range currentFieldIface {
			currentFieldStr[fmt.Sprint(k)] = v
		}
	}

	for k, v := range currentFieldStr {
		currentFieldStr[strings.ToLower(k)] = v
	}

	if len(path) == 1 {
		val, ok := currentFieldStr[firstInPath]
		if !ok {
			return "", false
		}
		return fmt.Sprint(val), true
	}

	return findValStrByPath(currentFieldStr[firstInPath], path[1:])
}
