package auth

import (
	"time"

	"github.com/labstack/echo"
)

type (
	// AccessTokenResponse defines the access token projection returned by API methods.
	AccessTokenResponse struct {
		TokenType    string `json:"token_type"`
		AccessToken  string `json:"access_token"`
		ExpiresIn    int64  `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}

	// AccessTokenSerializer defines functionality for serializing access tokens.
	AccessTokenSerializer struct {
		C echo.Context
		Token
		AccessToken  string
		RefreshToken string
	}
)

// Response returns the API response for a given access token.
func (s *AccessTokenSerializer) Response() AccessTokenResponse {
	return AccessTokenResponse{
		TokenType:    "bearer",
		AccessToken:  s.AccessToken,
		ExpiresIn:    int64(time.Until(s.Token.Expires) / time.Second),
		RefreshToken: s.RefreshToken,
	}
}
