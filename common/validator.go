package common

import "github.com/go-playground/validator"

// Validator definess a struct for validating objects.
type Validator struct {
	Validator *validator.Validate
}

// Validate validates the given object based on its constraints.
func (cv Validator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}
