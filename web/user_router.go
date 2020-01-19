package web

import (
	"net/http"

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

	// usersGroup := r.Group("/users", auth.RequireAdminRole())
	// usersGroup.GET("", r.findUsers)
	// usersGroup.POST("", r.createUser)
	// usersGroup.GET("/:id", r.getUser)
	// usersGroup.PUT("/:id", r.updateUser)
	// usersGroup.DELETE("/:id", r.deleteUser)
}

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

	user, err = r.userService.UpdateUser(&validator.user)
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

// func (r *userRouter) getUser(c echo.Context) error {
// 	id, err := bindUserID(c)
// 	if err != nil {
// 		return err
// 	}

// 	user, err := GetUserByID(uint(id))
// 	if err != nil {
// 		return err
// 	}

// 	serializer := UserSerializer{c, user}
// 	return c.JSON(http.StatusOK, serializer.Response())
// }

// func (r *userRouter) createUser(ec echo.Context) error {
// 	c, _ := ec.(*context)
// 	validator := NewUserModelValidator()
// 	if err := validator.Bind(c); err != nil {
// 		return err
// 	}

// 	if err := ensureUsernameIsNotTaken(validator.Username); err != nil {
// 		return err
// 	}

// 	if err := validator.userModel.Create(); err != nil {
// 		return err
// 	}

// 	serializer := UserSerializer{c, &validator.userModel}
// 	return c.JSON(http.StatusCreated, serializer.Response())
// }

// func (r *userRouter) findUsers(ec echo.Context) error {
// 	c, _ := ec.(*context)
// 	pr := c.BindPaging()
// 	users, totalCount, err := Find(pr)

// 	if err != nil {
// 		return err
// 	}

// 	serializer := UserListSerializer{c, users}
// 	return c.Page(http.StatusOK, pr, totalCount, serializer.Response())
// }

// func (r *userRouter) updateUser(ec echo.Context) error {
// 	c, _ := ec.(*context)
// 	id, err := bindUserID(c)
// 	if err != nil {
// 		return err
// 	}

// 	user, err := GetUserByID(uint(id))
// 	if err != nil {
// 		return err
// 	}

// 	validator := NewUserModelValidatorFillWith(*user)
// 	if err := validator.Bind(c); err != nil {
// 		return err
// 	}

// 	if validator.Username != validator.userModel.Username {
// 		if err := ensureUsernameIsNotTaken(validator.Username); err != nil {
// 			return err
// 		}
// 	}

// 	if err := validator.userModel.Save(); err != nil {
// 		return err
// 	}

// 	serializer := UserSerializer{c, &validator.userModel}
// 	return c.JSON(http.StatusOK, serializer.Response())
// }

// func (r *userRouter) deleteUser(c echo.Context) error {
// 	id, err := bindUserID(c)
// 	if err != nil {
// 		return err
// 	}

// 	user, err := GetUserByID(uint(id))
// 	if err != nil {
// 		return err
// 	}

// 	if err := user.Delete(); err != nil {
// 		return err
// 	}

// 	return c.NoContent(http.StatusNoContent)
// }

// func (r *userRouter) ensureUsernameIsNotTaken(username string) error {
// 	_, err := GetUserByUsername(username)

// 	if err == nil {
// 		err := errors.Conflict.Newf("Username '%s' already taken", username)
// 		return errors.AddContext(err, "username", "unique")
// 	}

// 	if errors.GetType(err) != errors.NotFound {
// 		return err
// 	}

// 	return nil
// }

// func (r *userRouter) bindUserID(c echo.Context) (uint, error) {
// 	id, err := strconv.ParseUint(c.Param("id"), 10, 32)

// 	if err != nil || id <= 0 {
// 		return 0, errors.BadRequest.New("User ID has to be a positive integer")
// 	}

// 	return uint(id), nil
// }
