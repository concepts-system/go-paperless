package application

import (
	"time"

	"github.com/google/uuid"

	"github.com/dgrijalva/jwt-go"
)

const (
	// TokenClaimSubject defines the token claim holding the token's subject.
	TokenClaimSubject = "sub"
	// TokenClaimUserID defines the token claim holding the user's ID.
	TokenClaimUserID = "user_id"
	// TokenClaimRoles defines the token claim holding the user's roles.
	TokenClaimRoles = "roles"
)

type (
	// Token defines a struct for holding authorization information.
	Token struct {
		Username       string
		UserID         uint
		Roles          []string
		Expires        time.Time
		RefreshExpires time.Time
	}

	// AccessTokenClaims defines all JWT (standard and custom) claims contained in an accesss tokens.
	AccessTokenClaims struct {
		jwt.StandardClaims
		UserID uint     `json:"user_id"`
		Roles  []string `json:"roles"`
	}

	// RefreshTokenClaims defines all JWT claims contained in a refresh token.
	RefreshTokenClaims struct {
		jwt.StandardClaims
		UserID uint `json:"user_id"`
	}
)

// GetAccessTokenClaims returns the JWT accesss token claims for the given Token instance.
func (t *Token) GetAccessTokenClaims(issuer, audience string) AccessTokenClaims {
	now := time.Now()
	id, _ := uuid.NewRandom()

	return AccessTokenClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        id.String(),
			Issuer:    issuer,
			Audience:  audience,
			Subject:   t.Username,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiresAt: t.Expires.Unix(),
		},

		UserID: t.UserID,
		Roles:  t.Roles,
	}
}

// GetRefreshTokenClaims returns the JWT refresh token claims for the given Token instance.
func (t *Token) GetRefreshTokenClaims(userID uint, issuer, audience string) RefreshTokenClaims {
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

		UserID: userID,
	}
}

// GrantsRole returns a boolean value indicating whether the token instance grants the given role.
func (t *Token) GrantsRole(role string) bool {
	for _, r := range t.Roles {
		if r == role {
			return true
		}
	}

	return false
}
