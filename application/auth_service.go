package application

import (
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/concepts-system/go-paperless/config"
	"github.com/concepts-system/go-paperless/domain"
	"github.com/concepts-system/go-paperless/errors"
)

const (
	// RoleUser defines the role every user belong to.
	RoleUser = "user"
	// RoleAdmin defines the role only admin users belong to.
	RoleAdmin = "admin"

	// TokenScopeAPI defines the scope for granting general API access.
	TokenScopeAPI = "api"
	// TokenScopeAuthRefresh defines the scope for granting refresh of
	// authentication.
	TokenScopeAuthRefresh = "auth:refresh"
)

// AuthService defines an application service for authentication and
// authorization use-cases.
type AuthService interface {
	// AuthenticateUserByCredentials tries to authenticate the user using the
	// given username and password and returns a new access token in case the
	// credentials are valid.
	AuthenticateUserByCredentials(username, password string) (*Token, error)

	// AuthenicateUserByRefreshToken tries to authenticate the user using the
	// given refresh token and returns a new access token in case the
	// provided refresh token is valid.
	AuthenicateUserByRefreshToken(token string) (*Token, error)

	// SignAccessToken signs the given token and returns the access token
	// encoded as a JWT.
	SignAccessToken(token *Token) (string, error)

	// SignRefreshToken signs the given token and returns the refresh token
	// encoded as a JWT.
	SignRefreshToken(token *Token) (string, error)

	// ExtractScopes extracts the token scopes from the given set of claims.
	ExtractScopes(claims jwt.MapClaims) []string

	// ExtractUsername extracts the username from the given set of claims.
	ExtractUsername(claims jwt.MapClaims) *string

	// ExtractRoles extracts the user's roles from the given set of claims.
	ExtractRoles(claims jwt.MapClaims) []string
}

type authServiceImpl struct {
	config           *config.Configuration
	users            domain.Users
	tokenKeyResolver TokenKeyResolver
	passwordHelper   *passwordHelper
}

// NewAuthService returns an auth service based on the given user repository
// and configuration.
func NewAuthService(
	config *config.Configuration,
	users domain.Users,
	tokenKeyResolver TokenKeyResolver,
) AuthService {
	return &authServiceImpl{
		config:           config,
		users:            users,
		tokenKeyResolver: tokenKeyResolver,
		passwordHelper:   &passwordHelper{},
	}
}

func (s *authServiceImpl) AuthenticateUserByCredentials(username, password string) (*Token, error) {
	user, err := s.users.GetByUsername(domain.Name(username))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve users")
	}

	if user == nil || !user.IsActive {
		return nil, s.badCredentialsError()
	}

	if err := s.passwordHelper.checkUserPassword(user, password); err != nil {
		return nil, s.badCredentialsError()
	}

	return s.userToken(user), nil
}

func (s *authServiceImpl) AuthenicateUserByRefreshToken(token string) (*Token, error) {
	refreshToken, err := jwt.Parse(token, s.tokenKeyResolver)
	if err != nil {
		return nil, UnauthorizedError.Newf("Invalid refresh token: %v", err)
	}

	if claims, ok := refreshToken.Claims.(jwt.MapClaims); ok && refreshToken.Valid {
		if !s.claimsScope(claims, TokenScopeAuthRefresh) {
			return nil, UnauthorizedError.Newf("Invalid token: Missing scope '%s'", TokenScopeAuthRefresh)
		}

		username, ok := claims[TokenClaimSubject].(string)
		if !ok {
			return nil, UnauthorizedError.New("Invalid token: Invalid subject claim")
		}

		user, err := s.users.GetByUsername(domain.Name(username))
		if err != nil {
			return nil, errors.Wrap(err, "Failed to read users")
		}

		if user == nil || !user.IsActive {
			return nil, s.badCredentialsError()
		}

		token, err := s.issueTokenForUser(domain.Name(username))
		if err != nil {
			return nil, err
		}

		return token, nil
	}

	return nil, UnauthorizedError.Newf("Invalid or expired refresh token")
}

func (s *authServiceImpl) SignAccessToken(token *Token) (string, error) {
	host := s.config.Server.PublicURL
	jwtToken := jwt.NewWithClaims(
		s.getSigningMethodHMAC(),
		token.GetAccessTokenClaims(
			host,
			host,
			TokenScopeAPI,
		),
	)

	accessToken, err := jwtToken.SignedString(s.config.Security.JWTSecret)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to sign access token: %s", err.Error())
	}

	return accessToken, nil
}

func (s *authServiceImpl) SignRefreshToken(token *Token) (string, error) {
	host := s.config.Server.PublicURL
	jwtToken := jwt.NewWithClaims(
		s.getSigningMethodHMAC(),
		token.GetRefreshTokenClaims(
			host,
			host,
			TokenScopeAuthRefresh,
		),
	)

	refreshToken, err := jwtToken.SignedString(s.config.Security.JWTSecret)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to sign refresh token: %s", err.Error())
	}

	return refreshToken, nil
}

func (s *authServiceImpl) ExtractScopes(claims jwt.MapClaims) []string {
	rawScope, ok := claims[TokenClaimScopes].(string)
	if !ok {
		return nil
	}

	rawScopes := strings.Split(rawScope, " ")
	scopes := make([]string, len(rawScopes))
	for i, scope := range rawScopes {
		scopes[i] = strings.Trim(scope, " \t")
	}

	return scopes
}

func (s *authServiceImpl) ExtractUsername(claims jwt.MapClaims) *string {
	if id, ok := claims[TokenClaimSubject].(string); ok {
		return &id
	}

	return nil
}

func (s *authServiceImpl) ExtractRoles(claims jwt.MapClaims) []string {
	rawRoles, ok := claims[TokenClaimRoles].([]interface{})
	if !ok {
		return nil
	}

	roles := make([]string, len(rawRoles))
	for i, role := range rawRoles {
		var ok bool
		roles[i], ok = role.(string)

		if !ok {
			return nil
		}
	}

	return roles
}

/* Helper Methods */

func (s *authServiceImpl) claimsScope(claims jwt.MapClaims, scope string) bool {
	for _, claimedScope := range s.ExtractScopes(claims) {
		if claimedScope == scope {
			return true
		}
	}

	return false
}

func (s *authServiceImpl) issueTokenForUser(username domain.Name) (*Token, error) {
	user, err := s.users.GetByUsername(username)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve user")
	}

	if user == nil {
		return nil, NotFoundError.Newf("User %s does not exist", username)
	}

	return s.userToken(user), nil
}

func (s *authServiceImpl) userToken(user *domain.User) *Token {
	return &Token{
		Username:       string(user.Username),
		Roles:          s.unwrapRoles(user),
		Expires:        time.Now().Add(s.config.Security.JWTExpirationTime),
		RefreshExpires: time.Now().Add(s.config.Security.JWTRefreshTime),
	}
}

func (s *authServiceImpl) getSigningMethodHMAC() *jwt.SigningMethodHMAC {
	switch s.config.Security.JWTAlgorithm {
	case "HS384":
		return jwt.SigningMethodHS384
	case "HS512":
		return jwt.SigningMethodHS512
	default:
		return jwt.SigningMethodHS256
	}
}

func (s *authServiceImpl) unwrapRoles(user *domain.User) []string {
	roles := []string{RoleUser}

	if user.IsAdmin {
		roles = append(roles, RoleAdmin)
	}

	return roles
}

func (s *authServiceImpl) badCredentialsError() error {
	return UnauthorizedError.New("Bad credentials")
}
