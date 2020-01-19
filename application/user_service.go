package application

import (
	"github.com/concepts-system/go-paperless/domain"
	"github.com/concepts-system/go-paperless/errors"
)

// UserService defines an application service for managing users use-cases.
//
// @ApplicationService
type UserService interface {
	// GetUserByID returns the user with the given ID or an error in case
	// no such user exists.
	GetUserByID(userID uint) (*domain.User, error)

	// FindUsers finds and returns users with respect to the given page request.
	FindUsers(pr domain.PageRequest) ([]domain.User, int64, error)

	// Creates the given new user with the desired password as clear-text.
	CreateNewUser(user *domain.User, password string) (*domain.User, error)

	// Update user updates all possible field of the given user.
	UpdateUser(user *domain.User) (*domain.User, error)

	// UpdateUserPassword updates the password of the user with the given ID.
	UpdateUserPassword(userID uint, currentPassword, newPassword string) error
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

func (s *userServiceImpl) GetUserByID(userID uint) (*domain.User, error) {
	return s.expectUserWithIDExists(domain.Identifier(userID))
}

func (s *userServiceImpl) FindUsers(pr domain.PageRequest) ([]domain.User, int64, error) {
	users, count, err := s.users.Find(pr)

	if err != nil {
		return nil, -1, errors.Wrap(err, "Failed to retreive users")
	}

	return users, int64(count), nil
}

func (s *userServiceImpl) CreateNewUser(
	user *domain.User,
	password string,
) (*domain.User, error) {
	if err := s.expectUsernameNotAlreadyTaken(user.Username); err != nil {
		return nil, err
	}

	if err := s.passwordHelper.setUserPassword(user, password); err != nil {
		return nil, errors.Wrap(err, "Failed to set user password")
	}

	user, err := s.users.Add(user)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create user")
	}

	return user, nil
}

func (s *userServiceImpl) UpdateUser(user *domain.User) (*domain.User, error) {
	originalUser, err := s.expectUserWithIDExists(user.ID)
	if err != nil {
		return nil, err
	}

	if user.Username != originalUser.Username {
		if err := s.expectUsernameNotAlreadyTaken(user.Username); err != nil {
			return nil, err
		}
	}

	originalUser.Username = user.Username
	originalUser.Forename = user.Forename
	originalUser.Surname = user.Surname

	user, err = s.users.Save(user)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to update user")
	}

	return user, nil
}

func (s *userServiceImpl) UpdateUserPassword(userID uint, currentPassword, newPassword string) error {
	user, err := s.expectUserWithIDExists(domain.Identifier(userID))
	if err != nil {
		return err
	}

	if err := s.passwordHelper.checkUserPassword(user, currentPassword); err != nil {
		err = BadRequestError.New("Incorrect current password")
		return errors.AddContext(err, "currentPassword", "value")
	}

	s.passwordHelper.setUserPassword(user, newPassword)
	if _, err := s.users.Save(user); err != nil {
		return err
	}

	return nil
}

func (s *userServiceImpl) expectUserWithIDExists(userID domain.Identifier) (*domain.User, error) {
	user, err := s.users.GetByID(userID)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to retrieve user")
	}

	if user == nil {
		return nil, NotFoundError.Newf("User with ID %d does not exist", userID)
	}

	return user, nil
}

func (s *userServiceImpl) expectUserWithUsernameExists(username domain.Name) (*domain.User, error) {
	user, err := s.users.GetByUsername(username)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to retrieve user")
	}

	if user == nil {
		return nil, NotFoundError.Newf("User with username '%s' does not exist", username)
	}

	return user, nil
}

func (s *userServiceImpl) expectUsernameNotAlreadyTaken(username domain.Name) error {
	user, err := s.users.GetByUsername(username)
	if err != nil {
		return errors.Wrapf(err, "Failed to retrieve user")
	}

	if user != nil {
		err := ConflictError.Newf("Username '%s' already taken", username)
		return errors.AddContext(err, "username", "unique")
	}

	return nil
}
