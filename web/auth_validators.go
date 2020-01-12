package web

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

// Bind binds the request to the request model.
func (validator *AuthenticationRequestValidator) Bind(c Context) error {
	return c.BindAndValidate(validator)
}

// NewPasswordAuthenticationRequestValidator returns a new instance of the respective validator.
func NewPasswordAuthenticationRequestValidator() PasswordAuthenticationRequestValidator {
	return PasswordAuthenticationRequestValidator{}
}

// Bind binds the request to the request model.
func (validator *PasswordAuthenticationRequestValidator) Bind(c Context) error {
	return c.BindAndValidate(validator)
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

// Bind binds the request to the request model.
func (validator *RefreshTokenAuthenticationRequestValidator) Bind(c Context) error {
	return c.BindAndValidate(validator)
}
