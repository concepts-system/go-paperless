package auth

import (
	"github.com/concepts-system/go-paperless/api"
)

// AuthenticationRequestValidator defines the validation rules for a general authentication request.
type AuthenticationRequestValidator struct {
	GrantType string `form:"grant_type" query:"grant_type" validate:"required"`
}

// NewAuthenticationRequestValidator returns a new instance of the respective validator.
func NewAuthenticationRequestValidator() AuthenticationRequestValidator {
	return AuthenticationRequestValidator{}
}

// PasswordAuthenticationRequestValidator defines the validation rules for an authentication request
// with grant type 'password'.
type PasswordAuthenticationRequestValidator struct {
	Username string `form:"username" validate:"required"`
	Password string `form:"password" validate:"required"`
}

// Bind binds the API request to the request model.
func (v *AuthenticationRequestValidator) Bind(c api.Context) error {
	return c.BindAndValidate(v)
}

// NewPasswordAuthenticationRequestValidator returns a new instance of the respective validator.
func NewPasswordAuthenticationRequestValidator() PasswordAuthenticationRequestValidator {
	return PasswordAuthenticationRequestValidator{}
}

// Bind binds the API request to the request model.
func (v *PasswordAuthenticationRequestValidator) Bind(c api.Context) error {
	return c.BindAndValidate(v)
}

// RefreshTokenAuthenticationRequestValidator defines the validation rules for an authentication request
// with grant type 'refresh_token'.
type RefreshTokenAuthenticationRequestValidator struct {
	Request struct {
		RefreshToken string `form:"refresh_token" validate:"required"`
	}
}

// NewRefreshTokenAuthenticationRequestValidator returns a new instance of the respective validator.
func NewRefreshTokenAuthenticationRequestValidator() RefreshTokenAuthenticationRequestValidator {
	return RefreshTokenAuthenticationRequestValidator{}
}

// Bind binds the API request to the request model.
func (v *RefreshTokenAuthenticationRequestValidator) Bind(c api.Context) error {
	return c.BindAndValidate(v)
}
