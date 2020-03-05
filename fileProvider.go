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

	logf("unsupported file type: %s", fileName)
	return nil
}

func findValStrByPath(i interface{}, path []string) (string, bool) { // todo: tests
	if len(path) == 0 {
		return "", false
	}

	currentField, ok := i.(map[interface{}]interface{})
	if !ok {
		return "", false
	}

	firstInPath := strings.ToLower(path[0])

	if len(path) == 1 {
		return fmt.Sprint(currentField[firstInPath]), true
	}

	return findValStrByPath(currentField[firstInPath], path[1:])
}
