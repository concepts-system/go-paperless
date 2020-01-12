package infrastructure

import (
	"github.com/concepts-system/go-paperless/domain"
	"github.com/jinzhu/gorm"
)

type usersGormImpl struct {
	db *Database
}

type userModel struct {
	gorm.Model
	Username string `gorm:"size:32;not_null;unique;unique_index"`
	Password string `gorm:"size:60;not_null"`
	Surname  string `gorm:"size:32;not_null;index"`
	Forename string `gorm:"size:32;not_null;index"`
	IsAdmin  bool   `gorm:"not_null"`
	IsActive bool   `gorm:"not_null"`
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

func (u usersGormImpl) GetByID(id domain.Identifier) (*domain.User, error) {
	var user userModel
	err := u.db.Where("id = ?", uint(id)).First(&user).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return u.mapUserModelToDomainEntity(&user), nil
}

func (u usersGormImpl) GetByUsername(username domain.Name) (*domain.User, error) {
	var user userModel
	err := u.db.
		Where("username = ?", string(username)).
		First(&user).
		Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return u.mapUserModelToDomainEntity(&user), nil
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

	return u.GetByID(domain.Identifier(userModel.ID))
}

func (u usersGormImpl) Save(user *domain.User) (*domain.User, error) {
	userModel := u.mapDomainEntityToUserModel(user)
	err := u.db.Save(userModel).Error

	if err != nil {
		return nil, err
	}

	return u.GetByID(domain.Identifier(userModel.ID))
}

func (u usersGormImpl) Delete(user *domain.User) error {
	return u.db.Delete(u).Error
}

func (u *usersGormImpl) mapUserModelToDomainEntity(user *userModel) *domain.User {
	return &domain.User{
		ID:       domain.Identifier(user.ID),
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

	model.ID = uint(user.ID)
	return model
}
