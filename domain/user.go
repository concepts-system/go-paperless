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

// Roles returns the roles for the given user.
func (u User) Roles() []Role {
	roles := make([]Role, 0)

	if u.IsAdmin {
		roles = append(roles, RoleAdmin)
	}

	return roles
}
