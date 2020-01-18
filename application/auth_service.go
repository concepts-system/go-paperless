package application

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/concepts-system/go-paperless/config"
	"github.com/concepts-system/go-paperless/domain"
	"github.com/concepts-system/go-paperless/errors"
)

// AuthService defines an application service for authentication and
// authorization use-cases.
//
// @ApplicationService
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
		return nil, errors.Wrap(err, "Failed to read users")
	}

	if user == nil {
		return nil, s.badCredentialsError()
	}

	if err := s.passwordHelper.checkUserPassword(user, password); err != nil {
		return nil, s.badCredentialsError()
	}

	return s.userToken(*user), nil
}

func (s *authServiceImpl) AuthenicateUserByRefreshToken(token string) (*Token, error) {
	refreshToken, err := jwt.Parse(token, s.tokenKeyResolver)
	if err != nil {
		return nil, UnauthorizedError.Newf("Invalid refresh token: %v", err)
	}

	if claims, ok := refreshToken.Claims.(jwt.MapClaims); ok && refreshToken.Valid {
		userID, ok := claims[TokenClaimUserID].(float64)

		if !ok || userID < 0 {
			return nil, UnauthorizedError.New("Invalid refressh token: Invalid user ID claim")
		}

		token, err := s.issueTokenForUser(uint(userID))
		if err != nil {
			return nil, err
		}

		return token, nil
	}

	return nil, UnauthorizedError.Newf("Invalid or expired refresh token")
}

func (s *authServiceImpl) SignAccessToken(token *Token) (string, error) {
	host := s.config.GetPublicURL().Host
	jwtToken := jwt.NewWithClaims(
		s.getSigningMethodHMAC(),
		token.GetAccessTokenClaims(
			host,
			host,
		),
	)

	// Sign and get the encoded token as a string using the secret
	accessToken, err := jwtToken.SignedString(s.config.GetJWTKey())
	if err != nil {
		return "", errors.Wrapf(err, "Failed to sign access token: %s", err.Error())
	}

	return accessToken, nil
}

func (s *authServiceImpl) SignRefreshToken(token *Token) (string, error) {
	host := s.config.GetPublicURL().Host
	jwtToken := jwt.NewWithClaims(
		s.getSigningMethodHMAC(),
		token.GetRefreshTokenClaims(
			token.UserID,
			host,
			host,
		),
	)

	// Sign and get the encoded token as a string using the secret
	refreshToken, err := jwtToken.SignedString(s.config.GetJWTKey())
	if err != nil {
		return "", errors.Wrapf(err, "Failed to sign refresh token: %s", err.Error())
	}

	return refreshToken, nil
}

func (s *authServiceImpl) issueTokenForUser(userID uint) (*Token, error) {
	user, err := s.users.GetByID(domain.Identifier(userID))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve user")
	}

	if user == nil {
		return nil, NotFoundError.Newf("No user with ID %d found", userID)
	}

	return s.userToken(*user), nil
}

func (s *authServiceImpl) userToken(user domain.User) *Token {
	return &Token{
		UserID:         uint(user.ID),
		Username:       string(user.Username),
		Roles:          s.unwrapRoles(user.Roles()),
		Expires:        time.Now().Add(s.config.GetJWTExpirationTime()),
		RefreshExpires: time.Now().Add(s.config.GetJWTRefreshTime()),
	}
}

func (s *authServiceImpl) getSigningMethodHMAC() *jwt.SigningMethodHMAC {
	switch s.config.GetJWTAlgorithm() {
	case config.JWTAlgorithmHS384:
		return jwt.SigningMethodHS384
	case config.JWTAlgorithmHS512:
		return jwt.SigningMethodHS512
	default:
		return jwt.SigningMethodHS256
	}
}

func (s *authServiceImpl) unwrapRoles(userRoles []domain.Role) []string {
	roles := make([]string, len(userRoles))

	for i, role := range userRoles {
		roles[i] = string(role)
	}

	return roles
}

func (s *authServiceImpl) badCredentialsError() error {
	return UnauthorizedError.New("Bad credentials")
}
