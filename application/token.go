package application

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	// TokenClaimSubject defines the token claim holding the token's subject.
	TokenClaimSubject = "sub"
	// TokenClaimRoles defines the token claim holding the user's roles.
	TokenClaimRoles = "roles"
	// TokenClaimScopes defines the token claim holding the token's scopes.
	TokenClaimScopes = "scope"
)

type (
	// Token defines a struct for holding authorization information.
	Token struct {
		Username       string
		Roles          []string
		Expires        time.Time
		RefreshExpires time.Time
	}

	// AccessTokenClaims defines all JWT (standard and custom) claims contained in an accesss tokens.
	AccessTokenClaims struct {
		jwt.StandardClaims
		Scope string   `json:"scope"`
		Roles []string `json:"roles"`
	}

	// RefreshTokenClaims defines all JWT claims contained in a refresh token.
	RefreshTokenClaims struct {
		jwt.StandardClaims
		Scope string `json:"scope"`
	}
)

// GetAccessTokenClaims returns the JWT accesss token claims for the given Token instance.
func (t *Token) GetAccessTokenClaims(issuer, audience, scope string) AccessTokenClaims {
	now := time.Now()

	return AccessTokenClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    issuer,
			Audience:  audience,
			Subject:   t.Username,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiresAt: t.Expires.Unix(),
		},

		Scope: scope,
		Roles: t.Roles,
	}
}

// GetRefreshTokenClaims returns the JWT refresh token claims for the given Token instance.
func (t *Token) GetRefreshTokenClaims(issuer, audience, scope string) RefreshTokenClaims {
	now := time.Now()

	return RefreshTokenClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    issuer,
			Audience:  audience,
			Subject:   t.Username,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiresAt: t.RefreshExpires.Unix(),
		},

		Scope: scope,
	}
}

// GrantsGroupMembership returns a boolean value indicating whether the token
// instance grants the given role.
func (t *Token) GrantsGroupMembership(group string) bool {
	for _, r := range t.Roles {
		if r == group {
			return true
		}
	}

	return false
}
