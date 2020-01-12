package web

import "github.com/go-playground/validator"

// Validator definess a struct for validating objects.
type Validator struct {
	Validator *validator.Validate
}

// Validate validates the given object based on its constraints.
func (validator Validator) Validate(i interface{}) error {
	return validator.Validator.Struct(i)
}
