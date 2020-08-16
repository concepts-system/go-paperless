package web

import (
	"github.com/concepts-system/go-paperless/domain"
)

type userPasswordUpdateValidator struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=8,max=255"`
}

type currentUserUpdateValidator struct {
	Surname  string `json:"surname" validate:"required,max=32"`
	Forename string `json:"forename" validate:"required,max=32"`

	user domain.User
}

type userCreationValidator struct {
	Username string `json:"username" validate:"required,alphanum,min=4,max=32"`
	Password string `json:"password" validate:"required,min=8,max=255"`
	Surname  string `json:"surname" validate:"required,max=32"`
	Forename string `json:"forename" validate:"required,max=32"`
	IsAdmin  bool   `json:"isAdmin"`
	IsActive bool   `json:"isActive"`

	user domain.User
}

type userUpdateValidator struct {
	Surname  string `json:"surname" validate:"required,max=32"`
	Forename string `json:"forename" validate:"required,max=32"`
	IsAdmin  bool   `json:"isAdmin"`
	IsActive bool   `json:"isActive"`

	user domain.User
}

func newPasswordUpdateValidator() *userPasswordUpdateValidator {
	return &userPasswordUpdateValidator{}
}

func newCurrentUserpdateValidatorOf(user *domain.User) *currentUserUpdateValidator {
	return &currentUserUpdateValidator{
		user:     *user,
		Forename: string(user.Forename),
		Surname:  string(user.Surname),
	}
}

func newUserCreationValidator() *userCreationValidator {
	return &userCreationValidator{
		IsActive: true,
		IsAdmin:  false,
	}
}

func newUserUpdateValidatorOf(user *domain.User) *userUpdateValidator {
	return &userUpdateValidator{
		user:     *user,
		Forename: string(user.Forename),
		Surname:  string(user.Surname),
		IsActive: user.IsActive,
		IsAdmin:  user.IsAdmin,
	}
}

func (v *userPasswordUpdateValidator) Bind(c *context) error {
	if err := c.BindAndValidate(v); err != nil {
		return err
	}

	return nil
}

func (v *currentUserUpdateValidator) Bind(c *context) error {
	if err := c.BindAndValidate(v); err != nil {
		return err
	}

	v.user.Surname = domain.Name(v.Surname)
	v.user.Forename = domain.Name(v.Forename)

	return nil
}

func (v *userCreationValidator) Bind(c *context) error {
	if err := c.BindAndValidate(v); err != nil {
		return err
	}

	v.user.Username = domain.Name(v.Username)
	v.user.Password = domain.Password(v.Password)
	v.user.Surname = domain.Name(v.Surname)
	v.user.Forename = domain.Name(v.Forename)
	v.user.IsAdmin = v.IsAdmin
	v.user.IsActive = v.IsActive

	return nil
}

func (v *userUpdateValidator) Bind(c *context) error {
	if err := c.BindAndValidate(v); err != nil {
		return err
	}

	v.user.Surname = domain.Name(v.Surname)
	v.user.Forename = domain.Name(v.Forename)
	v.user.IsAdmin = v.IsAdmin
	v.user.IsActive = v.IsActive

	return nil
}
