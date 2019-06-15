package users

import (
	"github.com/concepts-system/go-paperless/api"
	"github.com/concepts-system/go-paperless/errors"
)

// PasswordUpdateValidator defines the validation rules for update password requests.
type PasswordUpdateValidator struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=8,max=255"`
}

// Bind binds the request body to a proper password update request.
func (v *PasswordUpdateValidator) Bind(c api.Context) error {
	if err := c.BindAndValidate(v); err != nil {
		return err
	}

	return nil
}

// NewPasswordUpdateValidator constructs a validator with default values.
func NewPasswordUpdateValidator() PasswordUpdateValidator {
	return PasswordUpdateValidator{}
}

// UserModelValidator defines the validation rules for user models.
type UserModelValidator struct {
	Username  string  `json:"username" validate:"required,alphanum,min=4,max=32"`
	Password  *string `json:"password"`
	Surname   string  `json:"surname" validate:"required,max=32"`
	Forename  string  `json:"forename" validate:"required,max=32"`
	IsAdmin   bool    `json:"isAdmin"`
	IsActive  bool    `json:"isActive"`
	userModel UserModel
}

// Bind binds the given request to a user model.
func (v *UserModelValidator) Bind(c api.Context) error {
	if err := c.BindAndValidate(v); err != nil {
		return err
	}

	if v.Password == nil {
		err := errors.BadRequest.New("Validation failed")
		return errors.AddContext(err, "password", "required")
	} else if *v.Password != "" && len(*v.Password) < 8 {
		err := errors.BadRequest.New("Validation failed")
		return errors.AddContext(err, "password", "min")
	} else if *v.Password != "" && len(*v.Password) > 255 {
		err := errors.BadRequest.New("Validation failed")
		return errors.AddContext(err, "password", "max")
	}

	v.userModel.Username = v.Username
	v.userModel.SetPassword(*v.Password)
	v.userModel.Surname = v.Surname
	v.userModel.Forename = v.Forename
	v.userModel.IsAdmin = v.IsAdmin
	v.userModel.IsActive = v.IsActive

	return nil
}

// NewUserModelValidator constructs a validator with default values.
func NewUserModelValidator() UserModelValidator {
	validator := UserModelValidator{}
	validator.IsActive = true
	validator.IsAdmin = false
	return validator
}

// NewUserModelValidatorFillWith constructs a validator with the values from the given user model.
func NewUserModelValidatorFillWith(userModel UserModel) UserModelValidator {
	password := ""
	validator := NewUserModelValidator()
	validator.userModel = userModel
	validator.Username = userModel.Username
	validator.Password = &password
	validator.Surname = userModel.Surname
	validator.Forename = userModel.Forename
	validator.IsAdmin = userModel.IsAdmin
	validator.IsActive = userModel.IsActive

	return validator
}
