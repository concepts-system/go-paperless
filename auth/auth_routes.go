package auth

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"

	"github.com/concepts-system/go-paperless/api"
	"github.com/concepts-system/go-paperless/common"
	"github.com/concepts-system/go-paperless/errors"
)

const (
	keyGrantType          = "grant_type"
	grantTypePassword     = "password"
	grantTypeRefreshToken = "refresh_token"
)

var (
	errorBadCredentials = errors.Unauthorized.New("Bad credentials")
)

// RegisterRoutes registers all related routes for managing users.
func RegisterRoutes(r *echo.Group) {
	authgroup := r.Group("/auth")
	authgroup.POST("/token", getToken)
}

func getToken(ec echo.Context) error {
	c, _ := ec.(api.Context)
	validator := NewAuthenticationRequestValidator()
	if err := validator.Bind(c); err != nil {
		return err
	}

	switch validator.GrantType {
	case grantTypePassword:
		return getAccessTokenByPassword(c)
	case grantTypeRefreshToken:
		return getAccessTokenByRefreshToken(c)
	default:
		return errors.BadRequest.New("Unknown grant type")
	}
}

func getAccessTokenByPassword(ec echo.Context) error {
	c, _ := ec.(api.Context)
	validator := NewPasswordAuthenticationRequestValidator()
	if err := validator.Bind(c); err != nil {
		return err
	}

	token, err := AuthenticateUserWithCredentials(validator.Username, validator.Password)
	if err != nil {
		return errorBadCredentials
	}

	return sendTokenResponse(c, token)
}

func getAccessTokenByRefreshToken(ec echo.Context) error {
	c, _ := ec.(api.Context)
	// Verify token
	validator := NewRefreshTokenAuthenticationRequestValidator()
	if err := validator.Bind(c); err != nil {
		return err
	}

	refreshToken, err := jwt.Parse(validator.Request.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return common.Config().GetJWTKey(), nil
	})

	if err != nil {
		return errors.Unauthorized.Wrap(err, "Invalid refresh token")
	}

	if claims, ok := refreshToken.Claims.(jwt.MapClaims); ok && refreshToken.Valid {
		userID, ok := claims[claimUserID].(float64)

		if !ok || userID < 0 {
			return errors.Unauthorized.New("Invalid refressh token: Invalid user ID claim")
		}

		token, err := AuthenticateUserByID(uint(userID))
		if err != nil {
			return err
		}

		return sendTokenResponse(c, token)
	}

	// Fall back to unauthenticated
	return errorUnauthorized
}

func sendTokenResponse(c echo.Context, token Token) error {
	accessToken, err := SignAccessToken(token)
	if err != nil {
		return err
	}

	refreshToken, err := SignRefreshToken(token)
	if err != nil {
		return err
	}

	serializer := AccessTokenSerializer{
		c,
		token,
		accessToken,
		refreshToken,
	}

	return c.JSON(http.StatusOK, serializer.Response())
}
