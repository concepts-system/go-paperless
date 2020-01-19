package web

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/concepts-system/go-paperless/application"
)

const (
	keyGrantType          = "grant_type"
	grantTypePassword     = "password"
	grantTypeRefreshToken = "refresh_token"
)

type authRouter struct {
	authService      application.AuthService
	tokenKeyResolver application.TokenKeyResolver
}

// NewAuthRouter creates a new router for auth functionality based on the given
// configuration.
func NewAuthRouter(
	authService application.AuthService,
	tokenKeyResolver application.TokenKeyResolver,
) Router {
	return &authRouter{
		authService:      authService,
		tokenKeyResolver: tokenKeyResolver,
	}
}

// DefineRoutes defines the routes for auth functionality.
func (r *authRouter) DefineRoutes(group *echo.Group, auth *AuthMiddleware) {
	authgroup := group.Group("/auth")
	authgroup.POST("/token", r.getToken)
}

func (r *authRouter) getToken(ec echo.Context) error {
	c := ec.(*context)

	validator := newAuthenticationRequestValidator()
	if err := validator.Bind(c); err != nil {
		return err
	}

	switch validator.GrantType {
	case grantTypePassword:
		return r.getAccessTokenByPassword(c)
	case grantTypeRefreshToken:
		return r.getAccessTokenByRefreshToken(c)
	default:
		return application.BadRequestError.New("Unknown grant type")
	}
}

func (r *authRouter) getAccessTokenByPassword(c *context) error {
	validator := newPasswordAuthenticationRequestValidator()
	if err := validator.Bind(c); err != nil {
		return err
	}

	token, err := r.authService.AuthenticateUserByCredentials(
		validator.Username,
		validator.Password,
	)

	if err != nil {
		return err
	}

	return r.sendTokenResponse(c, token)
}

func (r *authRouter) getAccessTokenByRefreshToken(c *context) error {
	validator := newRefreshTokenAuthenticationRequestValidator()
	if err := validator.Bind(c); err != nil {
		return err
	}

	token, err := r.authService.AuthenicateUserByRefreshToken(validator.Request.RefreshToken)
	if err != nil {
		return err
	}

	return r.sendTokenResponse(c, token)
}

func (r *authRouter) sendTokenResponse(c *context, token *application.Token) error {
	accessToken, err := r.authService.SignAccessToken(token)
	if err != nil {
		return err
	}

	refreshToken, err := r.authService.SignRefreshToken(token)
	if err != nil {
		return err
	}

	serializer := accessTokenSerializer{
		c,
		token,
		accessToken,
		refreshToken,
	}

	return c.JSON(http.StatusOK, serializer.Response())
}
