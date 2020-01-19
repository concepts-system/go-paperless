package web

import (
	"time"

	"github.com/concepts-system/go-paperless/application"
)

type (
	// AccessTokenResponse defines the access token projection returned by API methods.
	accessTokenResponse struct {
		TokenType    string `json:"token_type"`
		AccessToken  string `json:"access_token"`
		ExpiresIn    int64  `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}

	// AccessTokenSerializer defines functionality for serializing access tokens.
	accessTokenSerializer struct {
		C *context
		*application.Token
		AccessToken  string
		RefreshToken string
	}
)

// Response returns the response for a given access token.
func (s accessTokenSerializer) Response() accessTokenResponse {
	return accessTokenResponse{
		TokenType:    "bearer",
		AccessToken:  s.AccessToken,
		ExpiresIn:    int64(time.Until(s.Token.Expires) / time.Second),
		RefreshToken: s.RefreshToken,
	}
}
