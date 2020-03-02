package configuration

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
)

func NewFileProvider(obj interface{}, fileName string) fileProvider {
	file, err := os.Open(fileName)
	if err != nil {
		return fileProvider{}
	}
	defer file.Close()

	if b, err := ioutil.ReadAll(file); err == nil {
		if fn := decodeFunc(fileName); fn != nil {
			err := fn(b, obj)
			log.Println(err)
		}
	}

	return fileProvider{}
}

type fileProvider struct{}

func (fileProvider) Provide(field reflect.StructField, v reflect.Value) bool {
	log.Println("str: ", v.Interface())

	//logf("fileProvider: set [%s] to field [%s] with tags [%v]", key, field.Name, field.Tag)
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
	return nil
}
