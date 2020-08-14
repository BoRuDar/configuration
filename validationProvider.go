package configuration

import (
	"fmt"
	"reflect"
	"errors"

	"github.com/go-playground/validator/v10"
)

var ErrByProvider = errors.New("provider error")

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
	if err := vP.provider.Provide(field, v, currentPath...); err != nil {
		return fmt.Errorf("validationProvider: %w, cannot be correct value: %v", ErrByProvider, err)
	}

	valStr := getValidateTag(field)
	validate := validator.New()
	return validate.Var(v.Interface(), valStr)
}