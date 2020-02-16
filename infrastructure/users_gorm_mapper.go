package infrastructure

import "github.com/concepts-system/go-paperless/domain"

type usersGormMapper struct{}

func newUsersGormMapper() *usersGormMapper {
	return &usersGormMapper{}
}

// MapUserModelToDomainEntity maps the given user model to the corresponding
// domain entity.
func (m *usersGormMapper) MapUserModelToDomainEntity(user *userModel) *domain.User {
	if user == nil {
		return nil
	}

	return &domain.User{
		Username: domain.Name(user.Username),
		Surname:  domain.Name(user.Surname),
		Forename: domain.Name(user.Forename),
		Password: domain.Password(user.Password),
		IsAdmin:  user.IsAdmin,
		IsActive: user.IsActive,
	}
}

// MapUserModelsToDomainEntities maps the given list of user models to a list
// containing the corresponding domain entities.
func (m *usersGormMapper) MapUserModelsToDomainEntities(users []userModel) []domain.User {
	if users == nil {
		return nil
	}

	domainEntities := make([]domain.User, len(users))

	for i, user := range users {
		domainEntities[i] = *m.MapUserModelToDomainEntity(&user)
	}

	return domainEntities
}

// MapDomainEntityToUserModel maps the given domain entity to the corresponding
// user model.
func (m *usersGormMapper) MapDomainEntityToUserModel(user *domain.User) *userModel {
	if user == nil {
		return nil
	}

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
