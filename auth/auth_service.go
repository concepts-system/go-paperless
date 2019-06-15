package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"github.com/concepts-system/go-paperless/common"
	"github.com/concepts-system/go-paperless/database"
	"github.com/concepts-system/go-paperless/errors"
)

const (
	// RoleAdmin defines the constant name of the admin role.
	RoleAdmin = "ADMIN"
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
func (t Token) GetAccessTokenClaims(issuer, audience string) AccessTokenClaims {
	now := time.Now()

	return AccessTokenClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        common.RandomString(32),
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
func (t Token) GetRefreshTokenClaims(userID uint, issuer, audience string) RefreshTokenClaims {
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
func (t Token) GrantsRole(role string) bool {
	for _, r := range t.Roles {
		if r == role {
			return true
		}
	}

	return false
}

type userDetails struct {
	ID       uint
	Username string
	Password string
	IsAdmin  bool
}

func (u userDetails) roles() []string {
	roles := make([]string, 0)

	if u.IsAdmin {
		roles = append(roles, RoleAdmin)
	}

	return roles
}

// AuthenticateUserByID returns a new access token for the given user ID if it exists.
// Note: This method does not verify anything and grants authorization for the given user ID
//       so be careful when using the result this method.
func AuthenticateUserByID(userID uint) (Token, error) {
	userDetails := findUserByID(userID)
	if userDetails == nil {
		return Token{}, errorBadCredentials
	}

	return getTokenForUser(*userDetails), nil
}

// AuthenticateUserWithCredentials tries to authenticate the user with the given username and password
// and returns a new access token in case the credentials are valid.
func AuthenticateUserWithCredentials(username, password string) (Token, error) {
	userDetails := findUserByUsername(username)
	if userDetails == nil {
		return Token{}, errorBadCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userDetails.Password), []byte(password)); err != nil {
		return Token{}, errorBadCredentials
	}

	return getTokenForUser(*userDetails), nil
}

// SignAccessToken signs the given token and returns the access token encoded as a JWT.
func SignAccessToken(token Token) (string, error) {
	host := common.Config().GetPublicURL().Host
	jwtToken := jwt.NewWithClaims(getSigningMethodHMAC(), token.GetAccessTokenClaims(
		host,
		host,
	))

	// Sign and get the encoded token as a string using the secret
	accessToken, err := jwtToken.SignedString(common.Config().GetJWTKey())
	if err != nil {
		return "", errors.Unexpected.Wrapf(err, "Failed to sign access token: %s", err.Error())
	}

	return accessToken, nil
}

// SignRefreshToken signs the given token and returns the refresh token encoded as a JWT.
func SignRefreshToken(token Token) (string, error) {
	host := common.Config().GetPublicURL().Host
	jwtToken := jwt.NewWithClaims(getSigningMethodHMAC(), token.GetRefreshTokenClaims(
		token.UserID,
		host,
		host,
	))

	// Sign and get the encoded token as a string using the secret
	refreshToken, err := jwtToken.SignedString(common.Config().GetJWTKey())
	if err != nil {
		return "", errors.Unexpected.Wrapf(err, "Failed to sign refresh token: %s", err.Error())
	}

	return refreshToken, nil
}

func getTokenForUser(details userDetails) Token {
	return Token{
		UserID:         details.ID,
		Username:       details.Username,
		Roles:          details.roles(),
		Expires:        time.Now().Add(common.Config().GetJWTExpirationTime()),
		RefreshExpires: time.Now().Add(common.Config().GetJWTRefreshTime()),
	}
}

func getSigningMethodHMAC() *jwt.SigningMethodHMAC {
	switch common.Config().GetJWTAlgorithm() {
	case common.JWTAlgorithmHS384:
		return jwt.SigningMethodHS384
	case common.JWTAlgorithmHS512:
		return jwt.SigningMethodHS512
	default:
		return jwt.SigningMethodHS256
	}
}

func findUserByID(userID uint) *userDetails {
	var userDetails userDetails
	database.DB().Table("users").
		Where("id = ? and is_active = true", userID).
		Select("id, username, password, is_admin").
		Scan(&userDetails)

	if userDetails.Username == "" {
		return nil
	}

	return &userDetails
}

func findUserByUsername(username string) *userDetails {
	var userDetails userDetails
	database.DB().Table("users").
		Where("username = ? and is_active = true", username).
		Select("id, username, password, is_admin").
		Scan(&userDetails)

	if userDetails.Username == "" {
		return nil
	}

	return &userDetails
}
