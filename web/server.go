package web

import (
	"fmt"

	"github.com/go-playground/validator"
	"github.com/kpango/glg"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/concepts-system/go-paperless/application"
	"github.com/concepts-system/go-paperless/config"
)

// Server represents the main object responsible for running the REST interface.
type Server struct {
	echo           *echo.Echo
	config         *config.Configuration
	authMiddleware *AuthMiddleware
}

// NewServer constructs a new server instance considering the given
// configuration.
func NewServer(config *config.Configuration) *Server {
	glg.Info("Initializing server...")

	server := Server{
		echo:           echo.New(),
		config:         config,
		authMiddleware: NewAuthMiddleware(application.ConfigTokenKeyResolver(config)),
	}

	// Configure echo instance
	server.echo = echo.New()
	server.echo.Debug = config.IsDevelopment()
	server.echo.HideBanner = true
	server.echo.HidePort = true
	server.echo.HTTPErrorHandler = errorHandler
	server.echo.Validator = Validator{Validator: validator.New()}

	// Register middlewares
	server.echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339}\t${method}\t${uri} -> status=${status} [${latency_human}] | ${error}\n",
	}))
	server.echo.Use(middleware.Recover())
	server.echo.Use(ExtendedContext)

	return &server
}

// Register registeres the given router with all its defined routes with the
// given server instance.
func (server *Server) Register(routers ...Router) {
	for _, router := range routers {
		router.DefineRoutes(
			server.echo.Group(""),
			server.authMiddleware,
		)
	}
}

// Start runs the server in a blocking way.
func (server *Server) Start() error {
	endpoint := fmt.Sprintf(":%d", server.config.GetPort())
	glg.Successf("Accepting connection on %s", endpoint)
	return server.echo.Start(endpoint)
}
