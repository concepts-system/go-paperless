package users

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/concepts-system/go-paperless/api"
	"github.com/concepts-system/go-paperless/auth"
	"github.com/concepts-system/go-paperless/errors"
)

var errorUsernameAlreadyExists = errors.Conflict.Newf("A user with the given username does already exist")

// RegisterRoutes registers all related routes for managing users.
func RegisterRoutes(r *echo.Group) {
	userGroup := r.Group("/user", auth.RequireAuthorization())
	userGroup.GET("/me", getCurrentUser)
	userGroup.PUT("/me", updateCurrentUser)
	userGroup.PUT("/me/password", updateCurrentUsersPassword)

	usersGroup := r.Group("/users", auth.RequireAdminRole())
	usersGroup.GET("", findUsers)
	usersGroup.POST("", createUser)
	usersGroup.GET("/:id", getUser)
	usersGroup.PUT("/:id", updateUser)
	usersGroup.DELETE("/:id", deleteUser)
}

func getCurrentUser(ec echo.Context) error {
	c, _ := ec.(api.Context)
	user, err := GetUserByID(*c.UserID)

	if err != nil {
		return err
	}

	serializer := UserSerializer{c, user}
	return c.JSON(http.StatusOK, serializer.Response())
}

func updateCurrentUser(ec echo.Context) error {
	c, _ := ec.(api.Context)
	user, err := GetUserByID(*c.UserID)

	if err != nil {
		return err
	}

	validator := NewUserModelValidatorFillWith(*user)
	if err := validator.Bind(c); err != nil {
		return err
	}

	if validator.userModel.Username != user.Username {
		if err := ensureUsernameIsNotTaken(validator.Username); err != nil {
			return err
		}
	}

	validator.userModel.Update(map[string]string{
		"username": validator.userModel.Username,
		"forename": validator.userModel.Forename,
		"surname":  validator.userModel.Surname,
	})

	serializer := UserSerializer{c, &validator.userModel}
	return c.JSON(http.StatusOK, serializer.Response())
}

func updateCurrentUsersPassword(ec echo.Context) error {
	c, _ := ec.(api.Context)
	validator := NewPasswordUpdateValidator()
	if err := validator.Bind(c); err != nil {
		return err
	}

	user, err := GetUserByID(*c.UserID)
	if err != nil {
		return err
	}

	if err := user.CheckPassword(validator.CurrentPassword); err != nil {
		err = errors.BadRequest.New("Incorrect current password")
		return errors.AddContext(err, "currentPassword", "value")
	}

	user.SetPassword(validator.NewPassword)
	if err := user.Update(map[string]string{"password": user.Password}); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func getUser(c echo.Context) error {
	id, err := bindUserID(c)
	if err != nil {
		return err
	}

	user, err := GetUserByID(uint(id))
	if err != nil {
		return err
	}

	serializer := UserSerializer{c, user}
	return c.JSON(http.StatusOK, serializer.Response())
}

func createUser(ec echo.Context) error {
	c, _ := ec.(api.Context)
	validator := NewUserModelValidator()
	if err := validator.Bind(c); err != nil {
		return err
	}

	if err := ensureUsernameIsNotTaken(validator.Username); err != nil {
		return err
	}

	if err := validator.userModel.Create(); err != nil {
		return err
	}

	serializer := UserSerializer{c, &validator.userModel}
	return c.JSON(http.StatusCreated, serializer.Response())
}

func findUsers(ec echo.Context) error {
	c, _ := ec.(api.Context)
	pr := c.BindPaging()
	users, totalCount, err := Find(pr)

	if err != nil {
		return err
	}

	serializer := UserListSerializer{c, users}
	return c.Page(http.StatusOK, pr, totalCount, serializer.Response())
}

func updateUser(ec echo.Context) error {
	c, _ := ec.(api.Context)
	id, err := bindUserID(c)
	if err != nil {
		return err
	}

	user, err := GetUserByID(uint(id))
	if err != nil {
		return err
	}

	validator := NewUserModelValidatorFillWith(*user)
	if err := validator.Bind(c); err != nil {
		return err
	}

	if validator.Username != validator.userModel.Username {
		if err := ensureUsernameIsNotTaken(validator.Username); err != nil {
			return err
		}
	}

	if err := validator.userModel.Save(); err != nil {
		return err
	}

	serializer := UserSerializer{c, &validator.userModel}
	return c.JSON(http.StatusOK, serializer.Response())
}

func deleteUser(c echo.Context) error {
	id, err := bindUserID(c)
	if err != nil {
		return err
	}

	user, err := GetUserByID(uint(id))
	if err != nil {
		return err
	}

	if err := user.Delete(); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func ensureUsernameIsNotTaken(username string) error {
	_, err := GetUserByUsername(username)

	if err == nil {
		err := errors.Conflict.Newf("Username '%s' already taken", username)
		return errors.AddContext(err, "username", "unique")
	}

	if errors.GetType(err) != errors.NotFound {
		return err
	}

	return nil
}

func bindUserID(c echo.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)

	if err != nil || id <= 0 {
		return 0, errors.BadRequest.New("User ID has to be a positive integer")
	}

	return uint(id), nil
}
