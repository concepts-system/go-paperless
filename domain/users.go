package domain

// Users defines an interface for managing the collection of all users.
//
// @DomainRepository
type Users interface {
	// GetByID returns the user with the given user ID or nil in case no such
	// user exists.
	GetByID(id Identifier) (*User, error)

	// GetByUsername returns the user with the given username or nil in case
	// no such user exists.
	GetByUsername(username Name) (*User, error)

	// Find returns the set of users and total count
	// with respect to the given page request.
	Find(page PageRequest) ([]User, Count, error)

	// Add adds a new user.
	Add(user *User) (*User, error)

	// Save saves the given user.
	Save(user *User) (*User, error)

	// Delete deletes the given user.
	Delete(user *User) error
}
