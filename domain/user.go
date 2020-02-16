package domain

type (
	// Name represents the type for common names.
	Name string

	// Password represents the type for a password.
	Password string
)

// User represents a user within the system.
type User struct {
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
