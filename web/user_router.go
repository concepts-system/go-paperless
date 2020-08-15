package web

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/concepts-system/go-paperless/application"
	"github.com/concepts-system/go-paperless/domain"
	"github.com/concepts-system/go-paperless/errors"
)

var errorUsernameAlreadyExists = application.ConflictError.Newf("A user with the given username does already exist")

type userRouter struct {
	userService application.UserService
}

// NewUserRouter creates a new router for user management using the given user
// service.
func NewUserRouter(userService application.UserService) Router {
	return &userRouter{
		userService: userService,
	}
}

// DefineRoutes defines the routes for auth functionality.
func (r *userRouter) DefineRoutes(group *echo.Group, auth *AuthMiddleware) {
	apiGroup := group.Group("/api", auth.RequireScope(application.TokenScopeAPI))
	userGroup := apiGroup.Group("/user", auth.RequireAuthentication())
	userGroup.GET("/me", r.getCurrentUser)
	userGroup.PATCH("/me", r.updateCurrentUser)
	userGroup.PUT("/me/password", r.updateCurrentUsersPassword)

	usersGroup := apiGroup.Group("/users", auth.RequireAdminRole())
	usersGroup.GET("", r.getUsers)
	usersGroup.POST("", r.createUser)
	usersGroup.GET("/:username", r.getUser)
	usersGroup.PATCH("/:username", r.updateUser)
	usersGroup.DELETE("/:username", r.deleteUser)
}

/* Handlers */

func (r *userRouter) getCurrentUser(ec echo.Context) error {
	c, _ := ec.(*context)
	user, err := r.userService.GetUserByUsername(*c.Username)

	if err != nil {
		return err
	}

	serializer := userSerializer{c, user}
	return c.JSON(http.StatusOK, serializer.Response())
}

func (r *userRouter) updateCurrentUser(ec echo.Context) error {
	c, _ := ec.(*context)

	user, err := r.userService.GetUserByUsername(*c.Username)
	if err != nil {
		return err
	}

	validator := newUserValidatorOf(user, false)
	if err := validator.Bind(c); err != nil {
		return err
	}

	validator.user.Username = domain.Name(*c.Username)
	validator.user.IsActive = user.IsActive
	validator.user.IsAdmin = user.IsAdmin

	user, err = r.userService.UpdateUser(&validator.user, nil)
	if err != nil {
		return err
	}

	serializer := userSerializer{c, user}
	return c.JSON(http.StatusOK, serializer.Response())
}

func (r *userRouter) updateCurrentUsersPassword(ec echo.Context) error {
	c, _ := ec.(*context)

	validator := newPasswordUpdateValidator()
	if err := validator.Bind(c); err != nil {
		return err
	}

	err := r.userService.UpdateUserPassword(
		*c.Username,
		validator.CurrentPassword,
		validator.NewPassword,
	)

	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (r *userRouter) getUsers(ec echo.Context) error {
	c, _ := ec.(*context)
	pr := c.BindPaging()
	users, totalCount, err := r.userService.GetUsers(pr.ToDomainPageRequest())

	if err != nil {
		return err
	}

	serializer := userListSerializer{c, users}
	return c.Page(http.StatusOK, pr, totalCount, serializer.Response())
}

func (r *userRouter) getUser(c echo.Context) error {
	user, err := r.userService.GetUserByUsername(r.bindUsername(c))
	if err != nil {
		return err
	}

	serializer := userSerializer{c, user}
	return c.JSON(http.StatusOK, serializer.Response())
}

func (r *userRouter) createUser(ec echo.Context) error {
	c, _ := ec.(*context)

	validator := newUserValidator(true)
	if err := validator.Bind(c); err != nil {
		return err
	}

	user, err := r.userService.CreateNewUser(
		&validator.user,
		*validator.Password,
	)

	if err != nil {
		return err
	}

	serializer := userSerializer{c, user}
	return c.JSON(http.StatusCreated, serializer.Response())
}

func (r *userRouter) updateUser(ec echo.Context) error {
	c, _ := ec.(*context)

	username := r.bindUsername(c)
	user, err := r.userService.GetUserByUsername(username)
	if err != nil {
		return err
	}

	validator := newUserValidatorOf(user, false)
	if err := validator.Bind(c); err != nil {
		return err
	}

	if username == *c.Username {
		if user.IsActive != validator.user.IsActive {
			err := application.BadRequestError.New("User may not alter his own active state")
			return errors.AddContext(err, "isActive", "const")
		}

		if user.IsAdmin != validator.user.IsAdmin {
			err := application.BadRequestError.New("User may not alter his own privileges")
			return errors.AddContext(err, "isAdmin", "const")
		}
	}

	validator.user.Username = domain.Name(*c.Username)
	user, err = r.userService.UpdateUser(&validator.user, validator.Password)
	if err != nil {
		return err
	}

	serializer := userSerializer{c, user}
	return c.JSON(http.StatusOK, serializer.Response())
}

func (r *userRouter) deleteUser(ec echo.Context) error {
	c, _ := ec.(*context)
	username := r.bindUsername(c)

	if username == *c.Username {
		return application.BadRequestError.New("User may not delete himself")
	}

	if err := r.userService.DeleteUser(username); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

/* Helper Methods */

func (r *userRouter) bindUsername(c echo.Context) string {
	return c.Param("username")
}
