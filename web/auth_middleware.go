package web

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/labstack/echo/v4"

	"github.com/concepts-system/go-paperless/application"
	"github.com/concepts-system/go-paperless/domain"
)

const (
	authorizationHeader        = "Authorization"
	authorizationBearerPrefix  = "Bearer "
	authorizationValueName     = "access_token"
	authorizationParameterName = "_token"
)

/* Token Extraction */

type tokenQueryParameterExtractor struct {
	ParameterName string
}

var authorizationHeaderExtractor = &request.PostExtractionFilter{
	Extractor: request.HeaderExtractor{authorizationHeader},

	Filter: func(headerValue string) (string, error) {
		length := len(authorizationBearerPrefix)

		if len(headerValue) >= length &&
			strings.ToLower(headerValue[0:length]) == strings.ToLower(authorizationBearerPrefix) {

			return headerValue[length:], nil
		}

		return headerValue, nil
	},
}

var tokenExtractor = &request.MultiExtractor{
	authorizationHeaderExtractor,
	request.ArgumentExtractor{authorizationValueName},
	&tokenQueryParameterExtractor{authorizationParameterName},
}

func (t *tokenQueryParameterExtractor) ExtractToken(r *http.Request) (string, error) {
	return r.URL.Query().Get(authorizationParameterName), nil
}

/* Middleware */

type (
	// filter defines a function type for defining a custom authorization condition.
	filter = func(userID uint, username string, roles []string) bool

	// AuthMiddleware defines which types of filters are provided by the auth middleware.
	AuthMiddleware struct {
		TokenKeyResolver application.TokenKeyResolver
	}
)

// NewAuthMiddleware creates a new auth middleware using the given
// token key resolver.
func NewAuthMiddleware(tokenKeyResolver application.TokenKeyResolver) *AuthMiddleware {
	return &AuthMiddleware{TokenKeyResolver: tokenKeyResolver}
}

// Require defines a middleware which checks for a valid authentication token
// and uses to custom function to authorize the token.
func (auth *AuthMiddleware) Require(filter filter) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ec echo.Context) error {
			c, _ := ec.(*context)
			clearAuthContext(c)

			// Verify token
			token, err := request.ParseFromRequest(
				c.Request(),
				tokenExtractor,
				auth.TokenKeyResolver,
			)

			// Token is invalid
			if err != nil {
				return application.UnauthorizedError.Newf("Invalid access token: %s", err.Error())
			}

			// Extract and verify claims
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				// Verify user ID claim
				id, ok := claims[application.TokenClaimUserID].(float64)

				if !ok || id <= 0 {
					return errorInvalidToken()
				}

				userID := uint(id)

				// Verify username claim
				username, ok := claims[application.TokenClaimSubject].(string)

				if !ok {
					return errorInvalidToken()
				}

				// Verify roles claim
				var rawClaims []interface{}
				rawClaims, ok = claims[application.TokenClaimRoles].([]interface{})
				claimedRoles := unwrapClaims(rawClaims)

				if claimedRoles == nil {
					return errorInsufficientPermission()
				}

				// Verify authorization condition
				if authorized := filter(userID, username, claimedRoles); !authorized {
					return errorInsufficientPermission()
				}

				// Fill auth context and continue to next handler
				setAuthContext(c, userID, username, claimedRoles)
				return h(c)
			}

			// Fallback to unauthorized
			return errorUnauthorized()
		}
	}
}

// RequireAuthorization returns a middleware handler function for protecting
// end-points needing user authentication.
// Any valid user ID claim will pass the middleware.
func (auth *AuthMiddleware) RequireAuthorization() echo.MiddlewareFunc {
	return auth.Require(func(_ uint, username string, _ []string) bool {
		return strings.TrimSpace(username) != ""
	})
}

// RequireRoles returns a middleware handler function for protecting end-points
// needing user authentication. The authenticated user needs also the given set
// of roles to get access granted.
func (auth *AuthMiddleware) RequireRoles(requiredRoles ...string) echo.MiddlewareFunc {
	return auth.Require(func(_ uint, _ string, roles []string) bool {
		if len(requiredRoles) == 0 {
			return true
		}

		// Check if each required role is claimed
		for requiredRole := range requiredRoles {
			hasRequiredRole := false

			for claimedRole := range roles {
				if requiredRole == claimedRole {
					hasRequiredRole = true
				}
			}

			if !hasRequiredRole {
				return false
			}
		}

		return true
	})
}

// RequireAdminRole returns a middleware handler function for protecting
// end-points needing user authentication.
// Only users having with 'ADMIN' role are allowed to access the end-point.
func (auth *AuthMiddleware) RequireAdminRole() echo.MiddlewareFunc {
	return auth.RequireRoles(string(domain.RoleAdmin))
}

func unwrapClaims(rawClaims []interface{}) []string {
	if rawClaims == nil {
		return nil
	}

	claims := make([]string, len(rawClaims))
	for i, claim := range rawClaims {
		var ok bool
		claims[i], ok = claim.(string)

		if !ok {
			return nil
		}
	}

	return claims
}

func clearAuthContext(c *context) {
	c.UserID = nil
	c.Username = nil
	c.Roles = nil
}

func setAuthContext(c *context, userID uint, username string, roles []string) {
	c.UserID = &userID
	c.Username = &username
	c.Roles = roles
}

/* Common Errors */

func errorInvalidToken() error {
	return application.UnauthorizedError.New("Invalid token")
}

func errorUnauthorized() error {
	return application.UnauthorizedError.New("Unauthorized")
}

func errorInsufficientPermission() error {
	return application.UnauthorizedError.New("Insufficient permissions")
}
