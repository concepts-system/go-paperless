package web

import "github.com/labstack/echo/v4"

// Router defines an abstraction for a REST router.
type Router interface {
	DefineRoutes(group *echo.Group, auth *AuthMiddleware)
}
