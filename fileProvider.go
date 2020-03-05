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
	log.Println("val: ", v.Interface())
	log.Printf("field: %+v\n", v.Interface())
	log.Println("path: ", path)

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
