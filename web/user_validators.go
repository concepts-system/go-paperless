package web

import (
	"github.com/concepts-system/go-paperless/application"
	"github.com/concepts-system/go-paperless/domain"
	"github.com/concepts-system/go-paperless/errors"
)

// userPasswordUpdateValidator defines the validation rules for update password requests.
type userPasswordUpdateValidator struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=8,max=255"`
}

// Bind binds the request body to a proper password update request.
func (v *userPasswordUpdateValidator) Bind(c *context) error {
	if err := c.BindAndValidate(v); err != nil {
		return err
	}

	return nil
}

// newPasswordUpdateValidator constructs a validator with default values.
func newPasswordUpdateValidator() *userPasswordUpdateValidator {
	return &userPasswordUpdateValidator{}
}

// userValidator defines the validation rules for user models.
type userValidator struct {
	passwordRequired bool

	Username string  `json:"username" validate:"required,alphanum,min=4,max=32"`
	Password *string `json:"password"`
	Surname  string  `json:"surname" validate:"required,max=32"`
	Forename string  `json:"forename" validate:"required,max=32"`
	IsAdmin  bool    `json:"isAdmin"`
	IsActive bool    `json:"isActive"`

	user domain.User
}

// Bind binds the given request to a user model.
func (v *userValidator) Bind(c *context) error {
	if err := c.BindAndValidate(v); err != nil {
		return err
	}

	if v.passwordRequired {
		if v.Password == nil {
			err := application.BadRequestError.New("Validation failed")
			return errors.AddContext(err, "password", "required")
		} else if *v.Password != "" && len(*v.Password) < 8 {
			err := application.BadRequestError.New("Validation failed")
			return errors.AddContext(err, "password", "min")
		} else if *v.Password != "" && len(*v.Password) > 255 {
			err := application.BadRequestError.New("Validation failed")
			return errors.AddContext(err, "password", "max")
		}
	}

	v.user.Username = domain.Name(v.Username)
	v.user.Surname = domain.Name(v.Surname)
	v.user.Forename = domain.Name(v.Forename)
	v.user.IsAdmin = v.IsAdmin
	v.user.IsActive = v.IsActive

	return nil
}

// newUserValidator constructs a validator with default values.
func newUserValidator(passwordRequired bool) *userValidator {
	return &userValidator{
		passwordRequired: passwordRequired,
		IsActive:         true,
		IsAdmin:          false,
	}
}

// newUserValidatorOf constructs a validator with the values from the given user model.
func newUserValidatorOf(user *domain.User, passwordRequired bool) *userValidator {
	password := ""
	validator := newUserValidator(passwordRequired)
	validator.user = *user
	validator.Username = string(user.Username)
	validator.Password = &password
	validator.Surname = string(user.Surname)
	validator.Forename = string(user.Forename)
	validator.IsAdmin = user.IsAdmin
	validator.IsActive = user.IsActive

	return validator
}
