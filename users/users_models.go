package users

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"

	"github.com/concepts-system/go-paperless/common"
	"github.com/concepts-system/go-paperless/database"
	"github.com/concepts-system/go-paperless/errors"
)

// UserModel defines the database model for users of the system.
type UserModel struct {
	gorm.Model
	Username string `gorm:"size:32;not_null;unique;unique_index"`
	Password string `gorm:"size:60;not_null"`
	Surname  string `gorm:"size:32;not_null;index"`
	Forename string `gorm:"size:32;not_null;index"`
	IsAdmin  bool   `gorm:"not_null"`
	IsActive bool   `gorm:"not_null"`
}

// TableName for UserModel entities.
func (UserModel) TableName() string {
	return "users"
}

// GetUserByID tries to find the user with the given ID.
func GetUserByID(id uint) (*UserModel, error) {
	var user UserModel
	err := database.DB().Where("id = ?", id).First(&user).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NotFound.Newf("User with ID '%d' not found", id)
		}

		return nil, errors.Wrap(err, "Failed to fetch user")
	}

	return &user, nil
}

// GetUserByUsername tries to find the user with the given username.
func GetUserByUsername(username string) (*UserModel, error) {
	var user UserModel
	err := database.DB().
		Where("username = ?", username).
		First(&user).
		Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NotFound.Newf("User with username '%s' not found", username)
		}

		return nil, errors.Wrap(err, "Failed to fetch user")
	}

	return &user, nil
}

// Find returns users within the system based on the given offset and limit params.
// The method will return the total count of found objects as 2nd parameter.
func Find(page common.PageRequest) ([]UserModel, int64, error) {
	var (
		users      []UserModel
		totalCount int64
	)

	err := database.DB().
		Order("surname, forename").
		Offset(page.Offset).
		Limit(page.Size).
		Find(&users).Error

	database.DB().Model(UserModel{}).Count(&totalCount)

	if err != nil {
		return nil, -1, err
	}

	return users, totalCount, nil
}

// Create persists the given UserModel instance in the database.
func (u *UserModel) Create() error {
	return database.DB().Create(u).Error
}

// Save saves (creates or updates) the given UserModel instance in the database.
func (u *UserModel) Save() error {
	return database.DB().Save(u).Error
}

// Update updates the given user model using the given expression.
func (u *UserModel) Update(expression interface{}) error {
	return database.DB().Model(u).Update(expression).Error
}

// Delete soft-removes the given UerModel instance from the database.
func (u *UserModel) Delete() error {
	return database.DB().Delete(u).Error
}

// CheckPassword checks whether the given password matches the user model's password hash.
func (u *UserModel) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// SetPassword sets the users password to the bcrypt hash of the given password.
func (u *UserModel) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return errors.Unexpected.Wrapf(err, "Failed to hash password")
	}

	u.Password = string(hash)
	return nil
}
