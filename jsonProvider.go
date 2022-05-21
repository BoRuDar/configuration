package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

const JSONFileProviderName = `JSONFileProvider`

var ErrFileMustHaveJSONExt = errors.New("file must have .json extension")

// NewJSONFileProvider creates new provider which reads values from JSON files.
func NewJSONFileProvider(fileName string) (fp *fileProvider) {
	return &fileProvider{fileName: fileName}
}

type fileProvider struct {
	fileName string
	fileData any
}

func (fileProvider) Name() string {
	return JSONFileProviderName
}

func (fp *fileProvider) Init(_ any) error {
	file, err := os.Open(fp.fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(strings.ToLower(fp.fileName), ".json") {
		return ErrFileMustHaveJSONExt
	}

	return json.Unmarshal(b, &fp.fileData)
}

func (fp fileProvider) Provide(field reflect.StructField, v reflect.Value) error {
	path := field.Tag.Get("file_json")
	if len(path) == 0 {
		// field doesn't have a proper tag
		return fmt.Errorf("%s: key is empty", JSONFileProviderName)
	}

	valStr, ok := findValStrByPath(fp.fileData, strings.Split(path, "."))
	if !ok {
		return fmt.Errorf("%s: findValStrByPath returns empty value", JSONFileProviderName)
	}

	return SetField(field, v, valStr)
}

func findValStrByPath(i any, path []string) (string, bool) {
	if len(path) == 0 {
		return "", false
	}
	firstInPath := strings.ToLower(path[0])

	currentFieldStr, ok := i.(map[string]any) // unmarshal from JSON
	if !ok {
		return "", false
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
