package configuration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

// NewFileProvider creates new provider which read values from files (json, yaml)
func NewFileProvider(fileName string) (fp fileProvider, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return fp, err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return fp, err
	}

	fn, err := decodeFunc(fileName)
	if err != nil {
		return fp, err
	}

	if err := fn(b, &fp.fileData); err != nil {
		return fp, err
	}
	return
}

type fileProvider struct {
	fileData interface{}
}

func (fp fileProvider) Provide(field reflect.StructField, v reflect.Value, path ...string) error {
	valStr, ok := findValStrByPath(fp.fileData, path)
	if !ok {
		return fmt.Errorf("fileProvider: findValStrByPath returns empty value")
	}

	return SetField(field, v, valStr)
}

func decodeFunc(fileName string) (func(data []byte, v interface{}) error, error) {
	fileName = strings.ToLower(fileName)

	if strings.HasSuffix(fileName, ".json") {
		return json.Unmarshal, nil
	}
	if strings.HasSuffix(fileName, ".yaml") {
		return yaml.Unmarshal, nil
	}
	if strings.HasSuffix(fileName, ".yml") {
		return yaml.Unmarshal, nil
	}

	return nil, fmt.Errorf("fileProvider: unsupported file type: %q", fileName)
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
