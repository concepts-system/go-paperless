package web

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/kpango/glg"
	"github.com/labstack/echo/v4"

	"github.com/concepts-system/go-paperless/application"
)

const (
	authorizationHeader        = echo.HeaderAuthorization
	authorizationBearerPrefix  = "Bearer "
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
	&tokenQueryParameterExtractor{authorizationParameterName},
}

func (t *tokenQueryParameterExtractor) ExtractToken(r *http.Request) (string, error) {
	return r.URL.Query().Get(authorizationParameterName), nil
}

/* Middleware */

type (
	// filter defines a function type for defining a custom authorization condition.
	filter = func(c *context) bool

	// AuthMiddleware defines which types of filters are provided by the auth middleware.
	AuthMiddleware struct {
		authService      application.AuthService
		tokenKeyResolver application.TokenKeyResolver
	}
)

// NewAuthMiddleware creates a new auth middleware using the given
// token key resolver and auth service.
func NewAuthMiddleware(
	authService application.AuthService,
	tokenKeyResolver application.TokenKeyResolver,
) *AuthMiddleware {
	return &AuthMiddleware{
		authService:      authService,
		tokenKeyResolver: tokenKeyResolver,
	}
}

func (auth *AuthMiddleware) authenticateRequest(c *context) error {
	glg.Debug("Authenticating request")
	auth.clearAuthContext(c)

	token, err := request.ParseFromRequest(
		c.Request(),
		tokenExtractor,
		auth.tokenKeyResolver,
	)

	if err != nil || !token.Valid {
		return application.UnauthorizedError.Newf(
			"Invalid token: %s",
			err.Error(),
		)
	}

	auth.extractClaims(c, token)
	return nil
}

// Require defines a middleware which checks for a valid authentication token
// and uses to custom function to authorize the token.
func (auth *AuthMiddleware) Require(filter filter) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ec echo.Context) error {
			c, _ := ec.(*context)

			if !c.IsAuthenticated() {
				auth.authenticateRequest(c)
			}

			if authorized := filter(c); !authorized {
				return auth.errorForbidden()
			}

			return h(c)
		}
	}
}

// RequireAuthentication returns a middleware handler function for protecting
// end-points needing user authentication.
//
// Any valid subject claim will pass the middleware.
func (auth *AuthMiddleware) RequireAuthentication() echo.MiddlewareFunc {
	return auth.Require(func(c *context) bool {
		glg.Debugf("Checking for user authentication")
		return c.IsAuthenticated()
	})
}

// RequireScope returns a middleware handler function for protecting end-points
// authentication. At least one of the given scopes needs to be granted in order
// to pass.
func (auth *AuthMiddleware) RequireScope(requiredScopes ...string) echo.MiddlewareFunc {
	return auth.Require(func(c *context) bool {
		glg.Debugf("Checking for scope: '%s'", strings.Join(requiredScopes, " "))

		for _, requiredScope := range requiredScopes {
			for _, claimedScope := range c.Scopes {
				if requiredScope == claimedScope {
					return true
				}
			}
		}

		return false
	})
}

// RequireRole returns a middleware handler function for protecting end-points
// needing user authentication. The authenticated user needs to have any role
// from the given set of roles in order to pass.
func (auth *AuthMiddleware) RequireRole(requiredRoles ...string) echo.MiddlewareFunc {
	return auth.Require(func(c *context) bool {
		glg.Debugf("Checking for any role of: %s", strings.Join(requiredRoles, ", "))

		for _, requiredRole := range requiredRoles {
			for _, claimedRole := range c.Roles {
				if requiredRole == claimedRole {
					return true
				}
			}
		}

		return false
	})
}

// RequireAdminRole returns a middleware handler function for protecting
// end-points needing user authentication.
//
// Only users having with 'admin' role are allowed to access the end-point.
func (auth *AuthMiddleware) RequireAdminRole() echo.MiddlewareFunc {
	return auth.RequireRole(application.RoleAdmin)
}

/* Helper Methods */

func (auth *AuthMiddleware) extractClaims(c *context, token *jwt.Token) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims == nil {
		return
	}

	c.Scopes = auth.authService.ExtractScopes(claims)
	c.Username = auth.authService.ExtractUsername(claims)
	c.Roles = auth.authService.ExtractRoles(claims)
}

func (auth *AuthMiddleware) clearAuthContext(c *context) {
	c.Scopes = nil
	c.Username = nil
	c.Roles = nil
}

/* Common Errors */

func (auth *AuthMiddleware) errorInvalidToken() error {
	return application.UnauthorizedError.New("Invalid token")
}

func (auth *AuthMiddleware) errorUnauthorized() error {
	return application.UnauthorizedError.New("Unauthorized")
}

func (auth *AuthMiddleware) errorForbidden() error {
	return application.ForbiddenError.New("Forbidden")
}
