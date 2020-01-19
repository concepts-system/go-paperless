package web

// authenticationRequestValidator defines the validation rules for a general authentication request.
type authenticationRequestValidator struct {
	GrantType string `form:"grant_type" query:"grant_type" validate:"required"`
}

// newAuthenticationRequestValidator returns a new instance of the respective validator.
func newAuthenticationRequestValidator() *authenticationRequestValidator {
	return &authenticationRequestValidator{}
}

// passwordAuthenticationRequestValidator defines the validation rules for an authentication request
// with grant type 'password'.
type passwordAuthenticationRequestValidator struct {
	Username string `form:"username" validate:"required"`
	Password string `form:"password" validate:"required"`
}

// Bind binds the request to the request model.
func (validator *authenticationRequestValidator) Bind(c *context) error {
	return c.BindAndValidate(validator)
}

// newPasswordAuthenticationRequestValidator returns a new instance of the respective validator.
func newPasswordAuthenticationRequestValidator() *passwordAuthenticationRequestValidator {
	return &passwordAuthenticationRequestValidator{}
}

// Bind binds the request to the request model.
func (validator *passwordAuthenticationRequestValidator) Bind(c *context) error {
	return c.BindAndValidate(validator)
}

// refreshTokenAuthenticationRequestValidator defines the validation rules for an authentication request
// with grant type 'refresh_token'.
type refreshTokenAuthenticationRequestValidator struct {
	Request struct {
		RefreshToken string `form:"refresh_token" validate:"required"`
	}
}

// newRefreshTokenAuthenticationRequestValidator returns a new instance of the respective validator.
func newRefreshTokenAuthenticationRequestValidator() *refreshTokenAuthenticationRequestValidator {
	return &refreshTokenAuthenticationRequestValidator{}
}

// Bind binds the request to the request model.
func (validator *refreshTokenAuthenticationRequestValidator) Bind(c *context) error {
	return c.BindAndValidate(validator)
}
