package application

import (
	"github.com/concepts-system/go-paperless/config"
	"github.com/dgrijalva/jwt-go"
)

type (
	// TokenKeyResolver defines a function type for a function that
	// obtains a verification key for a given token.
	TokenKeyResolver = func(token *jwt.Token) (interface{}, error)
)

// ConfigTokenKeyResolver returns a token key resolver using the key from
// the given config.
func ConfigTokenKeyResolver(config *config.Configuration) TokenKeyResolver {
	return func(token *jwt.Token) (interface{}, error) {
		b := config.GetJWTKey()
		return b, nil
	}
}
