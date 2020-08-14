package configuration

import (
	"github.com/go-playground/validator/v10"
	"reflect"
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
	vP.provider.Provide(field, v, currentPath...)

	valStr := getValidateTag(field)
	validate := validator.New()
	return validate.Var(v.Interface(), valStr)
}