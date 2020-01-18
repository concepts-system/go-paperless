package domain

// RoleAdmin defines the constant name of the admin role.
const RoleAdmin = Role("ADMIN")

// Name represents the type for common names.
//
// @ValueObject
type Name string

// Password represents the type for a password.
//
// @ValueObject
type Password string

// Role represents a role a user might have.
//
// @ValueObject
type Role string

// User represents a user within the system.
//
// @DomainEntity
type User struct {
	ID       Identifier
	Username Name
	Password Password
	Surname  Name
	Forename Name
	IsAdmin  bool
	IsActive bool
}

// NewUser creates a new, valid user based on the given values.
func NewUser(user User) *User {
	return &User{
		Username: user.Username,
		Surname:  user.Surname,
		Forename: user.Forename,
		IsAdmin:  user.IsAdmin,
		IsActive: user.IsActive,
	}
}

// Roles returns the roles for the given user.
func (u *User) Roles() []Role {
	roles := make([]Role, 0)

	if u.IsAdmin {
		roles = append(roles, RoleAdmin)
	}

	return roles
}
