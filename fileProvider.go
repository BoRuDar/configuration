package configuration

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"reflect"
)

func NewFileProvider(obj interface{}, file io.Reader) fileProvider {
	if b, err := ioutil.ReadAll(file); err == nil {
		log.Printf("file: %s\n", b)

		err = json.Unmarshal(b, obj)
		if err != nil {
			log.Fatalf("unmarshall: %v", err)
		}
	}

	return fileProvider{}
}

type fileProvider struct {
	r io.Reader
}

func (fileProvider) Provide(field reflect.StructField, v reflect.Value) bool {
	log.Println("str: ", v.Interface())

	//SetField(field, v, key)
	//logf("fileProvider: set [%s] to field [%s] with tags [%v]", key, field.Name, field.Tag)
	return true
}
