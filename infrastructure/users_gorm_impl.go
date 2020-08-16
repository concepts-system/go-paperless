package infrastructure

import (
	"time"

	"github.com/concepts-system/go-paperless/domain"
	"github.com/jinzhu/gorm"
)

type usersGormImpl struct {
	db     *Database
	mapper *usersGormMapper
}

type userModel struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
	Username  string     `gorm:"size:32;not_null;index:username"`
	Password  string     `gorm:"size:60;not_null"`
	Surname   string     `gorm:"size:32;not_null"`
	Forename  string     `gorm:"size:32;not_null"`
	IsAdmin   bool       `gorm:"not_null"`
	IsActive  bool       `gorm:"not_null"`
}

func (userModel) TableName() string {
	return "users"
}

// NewUsers creates a new users domain repository.
func NewUsers(db *Database) domain.Users {
	return usersGormImpl{
		db:     db,
		mapper: newUsersGormMapper(),
	}
}

func (u usersGormImpl) GetByUsername(username domain.Name) (*domain.User, error) {
	user, err := u.getUserModelByUsername(string(username))

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return u.mapper.MapUserModelToDomainEntity(user), nil
}

func (u usersGormImpl) Find(page domain.PageRequest) ([]domain.User, domain.Count, error) {
	var (
		users      []userModel
		totalCount int64
	)

	err := u.db.
		Order("surname, forename").
		Offset(page.Offset).
		Limit(page.Size).
		Find(&users).
		Count(&totalCount).
		Error

	if err != nil {
		return nil, -1, err
	}

	return u.mapper.MapUserModelsToDomainEntities(users), domain.Count(totalCount), nil
}

func (u usersGormImpl) Add(user *domain.User) (*domain.User, error) {
	userModel := u.mapper.MapDomainEntityToUserModel(user)
	err := u.db.Create(userModel).Scan(userModel).Error

	if err != nil {
		return nil, err
	}

	return u.mapper.MapUserModelToDomainEntity(userModel), nil
}

func (u usersGormImpl) Update(user *domain.User) (*domain.User, error) {
	model, err := u.getUserModelByUsername(string(user.Username))
	if err != nil {
		return nil, err
	}

	updatedModel := u.mapper.MapDomainEntityToUserModel(user)
	updatedModel.ID = model.ID

	if err := u.db.Save(updatedModel).Error; err != nil {
		return nil, err
	}

	return u.GetByUsername(user.Username)
}

func (u usersGormImpl) Delete(user *domain.User) error {
	model, err := u.getUserModelByUsername(string(user.Username))
	if model == nil {
		return err
	}

	return u.db.Delete(model).Error
}

func (u *usersGormImpl) getUserModelByUsername(username string) (*userModel, error) {
	var user userModel
	err := u.db.
		Where("username = ?", string(username)).
		First(&user).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
