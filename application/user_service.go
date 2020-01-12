package application

import (
	"github.com/concepts-system/go-paperless/domain"
	"github.com/concepts-system/go-paperless/errors"
)

// UserService defines an application service for managing users use-cases.
//
// @ApplicationService
type UserService interface {
	// FindUsers finds and returns users with respect to the given page request.
	FindUsers(pr domain.PageRequest) ([]domain.User, int64, error)

	// Creates the given new user with the desired password as clear-text.
	CreateNewUser(user *domain.User, password string) (*domain.User, error)
}

type userServiceImpl struct {
	users          domain.Users
	passwordHelper *passwordHelper
}

// NewUserService creates a new user service.
func NewUserService(users domain.Users) UserService {
	return &userServiceImpl{
		users:          users,
		passwordHelper: &passwordHelper{},
	}
}

func (s userServiceImpl) FindUsers(pr domain.PageRequest) ([]domain.User, int64, error) {
	users, count, err := s.users.Find(pr)

	if err != nil {
		return nil, -1, errors.Wrap(err, "Failed to retreive users")
	}

	return users, int64(count), nil
}

func (s userServiceImpl) CreateNewUser(
	user *domain.User,
	password string,
) (*domain.User, error) {
	existingUser, err := s.users.GetByUsername(user.Username)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to check for existing user")
	}

	if existingUser != nil {
		err := ConflictError.Newf("Username '%s' already taken", user.Username)
		return nil, errors.AddContext(err, "username", "unique")
	}

	if err := s.passwordHelper.setUserPassword(user, password); err != nil {
		return nil, errors.Wrap(err, "Failed to set user password")
	}

	user, err = s.users.Add(user)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create user")
	}

	return user, nil
}
