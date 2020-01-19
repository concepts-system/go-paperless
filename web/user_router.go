package web

import (
	"net/http"
	"strconv"

	"github.com/concepts-system/go-paperless/domain"
	"github.com/concepts-system/go-paperless/errors"

	"github.com/labstack/echo/v4"

	"github.com/concepts-system/go-paperless/application"
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
	apiGroup := group.Group("/api")
	userGroup := apiGroup.Group("/user", auth.RequireAuthorization())
	userGroup.GET("/me", r.getCurrentUser)
	userGroup.PUT("/me", r.updateCurrentUser)
	userGroup.PUT("/me/password", r.updateCurrentUsersPassword)

	usersGroup := apiGroup.Group("/users", auth.RequireAdminRole())
	usersGroup.GET("", r.findUsers)
	usersGroup.POST("", r.createUser)
	usersGroup.GET("/:id", r.getUser)
	usersGroup.PUT("/:id", r.updateUser)
	usersGroup.DELETE("/:id", r.deleteUser)
}

/* Handlers */

func (r *userRouter) getCurrentUser(ec echo.Context) error {
	c, _ := ec.(*context)
	user, err := r.userService.GetUserByID(*c.UserID)

	if err != nil {
		return err
	}

	serializer := userSerializer{c, user}
	return c.JSON(http.StatusOK, serializer.Response())
}

func (r *userRouter) updateCurrentUser(ec echo.Context) error {
	c, _ := ec.(*context)

	user, err := r.userService.GetUserByID(*c.UserID)
	if err != nil {
		return err
	}

	validator := newUserValidatorOf(user, false)
	if err := validator.Bind(c); err != nil {
		return err
	}

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
		*c.UserID,
		validator.CurrentPassword,
		validator.NewPassword,
	)

	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (r *userRouter) findUsers(ec echo.Context) error {
	c, _ := ec.(*context)
	pr := c.BindPaging()
	users, totalCount, err := r.userService.FindUsers(pr.ToDomainPageRequest())

	if err != nil {
		return err
	}

	serializer := userListSerializer{c, users}
	return c.Page(http.StatusOK, pr, totalCount, serializer.Response())
}

func (r *userRouter) getUser(c echo.Context) error {
	id, err := r.bindUserID(c)
	if err != nil {
		return err
	}

	user, err := r.userService.GetUserByID(uint(id))
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

	id, err := r.bindUserID(c)
	if err != nil {
		return err
	}

	user, err := r.userService.GetUserByID(uint(id))
	if err != nil {
		return err
	}

	validator := newUserValidatorOf(user, false)
	if err := validator.Bind(c); err != nil {
		return err
	}

	if r.isCurrentUserID(c, validator.user.ID) {
		if user.IsActive != validator.user.IsActive {
			err := application.BadRequestError.New("User may not alter his active state")
			return errors.AddContext(err, "isActive", "const")
		}

		if user.IsAdmin != validator.user.IsAdmin {
			err := application.BadRequestError.New("User may not alter his own privileges")
			return errors.AddContext(err, "isAdmin", "const")
		}
	}

	user, err = r.userService.UpdateUser(&validator.user, validator.Password)
	if err != nil {
		return err
	}

	serializer := userSerializer{c, user}
	return c.JSON(http.StatusOK, serializer.Response())
}

func (r *userRouter) deleteUser(ec echo.Context) error {
	c, _ := ec.(*context)
	id, err := r.bindUserID(c)
	if err != nil {
		return err
	}

	if r.isCurrentUserID(c, domain.Identifier(id)) {
		return application.BadRequestError.New("User may not delete himself")
	}

	if err = r.userService.DeleteUser(domain.Identifier(id)); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

/* Helper Methods */

func (r *userRouter) bindUserID(c echo.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)

	if err != nil || id <= 0 {
		return 0, application.BadRequestError.New("User ID has to be a positive integer")
	}

	return uint(id), nil
}

func (r *userRouter) isCurrentUserID(c *context, id domain.Identifier) bool {
	if *c.UserID == uint(id) {
		return true
	}

	return false
}
