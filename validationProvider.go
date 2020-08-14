package configuration

import (
	"log"
	"reflect"

	"github.com/go-playground/validator/v10"
)

// NewValidationProvider creates new provider which validates the provided value against the validate tag.
func NewValidationProvider(p Provider) validationProvider {
	return validationProvider{
		provider: p,
	}
}

type validationProvider struct{
	provider Provider
}

func (vP validationProvider) Provide(field reflect.StructField, v reflect.Value, currentPath ...string) error {
	err := vP.provider.Provide(field, v, currentPath...)
	if err != nil {
		log.Println(err)
	}

	valStr := getValidateTag(field)
	validate := validator.New()
	return validate.Var(v.Interface(), valStr)
}