package auth

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/labstack/echo"

	"github.com/concepts-system/go-paperless/api"
	"github.com/concepts-system/go-paperless/common"
	"github.com/concepts-system/go-paperless/errors"
)

const (
	authorizationHeader        = "Authorization"
	authorizationBearerPrefix  = "Bearer "
	authorizationValueName     = "access_token"
	authorizationParameterName = "_token"

	claimSubject = "sub"
	claimUserID  = "user_id"
	claimRoles   = "roles"
)

var (
	errorInvalidToken           = errors.Unauthorized.New("Invalid token")
	errorUnauthorized           = errors.Unauthorized.New("Unauthorized")
	errorInsufficientPermission = errors.Forbidden.New("Insufficient permissions")
)

// Filter defines a function type for defining a custom authorization condition.
type Filter = func(userID uint, username string, roles []string) bool

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

func clearAuthContext(c *api.Context) {
	c.UserID = nil
	c.Username = nil
	c.Roles = nil
}

func setAuthContext(c *api.Context, userID uint, username string, roles []string) {
	c.UserID = &userID
	c.Username = &username
	c.Roles = roles
}

// Middleware defines a middleware which checks for a valid authentication token and uses to custom function
// to authorize the token.
func Middleware(filter Filter) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ec echo.Context) error {
			c, _ := ec.(api.Context)
			clearAuthContext(&c)

			// Verify token
			token, err := request.ParseFromRequest(c.Request(), tokenExtractor, func(token *jwt.Token) (interface{}, error) {
				b := common.Config().GetJWTKey()
				return b, nil
			})

			// Token is invalid
			if err != nil {
				return errors.Unauthorized.Wrap(err, "Invalid access token")
			}

			// Extract and verify claims
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				// Verify user ID claim
				id, ok := claims[claimUserID].(float64)

				if !ok || id <= 0 {
					return errorInvalidToken
				}

				userID := uint(id)

				// Verify username claim
				username, ok := claims[claimSubject].(string)

				if !ok {
					return errorInvalidToken
				}

				// Verify roles claim
				var rawClaims []interface{}
				rawClaims, ok = claims[claimRoles].([]interface{})
				claimedRoles := unwrapClaims(rawClaims)

				if claimedRoles == nil {
					return errorInsufficientPermission
				}

				// Verify authorization condition
				if authorized := filter(userID, username, claimedRoles); !authorized {
					return errorInsufficientPermission
				}

				// Fill auth context and continue to next handler
				setAuthContext(&c, userID, username, claimedRoles)
				return h(c)
			}

			// Fallback to unauthorized
			return errorUnauthorized
		}
	}
}

// RequireAuthorization returns a middleware handler function for protecting end-points needing user authentication.
// Any valid user ID claim will pass the middleware.
func RequireAuthorization() echo.MiddlewareFunc {
	return Middleware(func(_ uint, username string, _ []string) bool {
		return strings.TrimSpace(username) != ""
	})
}

// RequireRoles returns a middleware handler function for protecting end-points needing user authentication.
// The authenticated user needs also the given set of roles to get access granted.
func RequireRoles(requiredRoles ...string) echo.MiddlewareFunc {
	return Middleware(func(_ uint, _ string, roles []string) bool {
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

// RequireAdminRole returns a middleware handler function for protecting end-points needing user authentication.
// Only users having with 'ADMIN' role are allowed to access the end-point.
func RequireAdminRole() echo.MiddlewareFunc {
	return RequireRoles(RoleAdmin)
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
