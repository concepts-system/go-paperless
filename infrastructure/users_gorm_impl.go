package infrastructure

import (
	"time"

	"github.com/concepts-system/go-paperless/domain"
	"github.com/jinzhu/gorm"
)

type usersGormImpl struct {
	db *Database
}

type userModel struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index:username"`
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
		db: db,
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

	return u.mapUserModelToDomainEntity(user), nil
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
		Find(&users).Error

	u.db.Model(userModel{}).Count(&totalCount)

	if err != nil {
		return nil, -1, err
	}

	return u.mapUserModelsToDomainEntities(users), domain.Count(totalCount), nil
}

func (u usersGormImpl) Add(user *domain.User) (*domain.User, error) {
	userModel := u.mapDomainEntityToUserModel(user)
	err := u.db.Create(userModel).Error

	if err != nil {
		return nil, err
	}

	return u.GetByUsername(user.Username)
}

func (u usersGormImpl) Update(user *domain.User) (*domain.User, error) {
	model, err := u.getUserModelByUsername(string(user.Username))
	if err != nil {
		return nil, err
	}

	updatedModel := u.mapDomainEntityToUserModel(user)
	updatedModel.ID = model.ID

	if err := u.db.Save(updatedModel).Error; err != nil {
		return nil, err
	}

	return u.GetByUsername(user.Username)
}

func (u usersGormImpl) Delete(user *domain.User) error {
	model, err := u.getUserModelByUsername(string(user.Username))
	if user == nil {
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

func (u *usersGormImpl) mapUserModelToDomainEntity(user *userModel) *domain.User {
	return &domain.User{
		Username: domain.Name(user.Username),
		Surname:  domain.Name(user.Surname),
		Forename: domain.Name(user.Forename),
		Password: domain.Password(user.Password),
		IsAdmin:  user.IsAdmin,
		IsActive: user.IsActive,
	}
}

func (u *usersGormImpl) mapUserModelsToDomainEntities(users []userModel) []domain.User {
	domainEntities := make([]domain.User, len(users))

	for i, user := range users {
		domainEntities[i] = *u.mapUserModelToDomainEntity(&user)
	}

	return domainEntities
}

func (u *usersGormImpl) mapDomainEntityToUserModel(user *domain.User) *userModel {
	model := &userModel{
		Username: string(user.Username),
		Surname:  string(user.Surname),
		Forename: string(user.Forename),
		Password: string(user.Password),
		IsAdmin:  user.IsAdmin,
		IsActive: user.IsActive,
	}

	return model
}
