package domain

// Users defines an interface for managing the collection of all users.
//
// @DomainRepository
type Users interface {
	// GetByUsername returns the user with the given username or nil in case
	// no such user exists.
	GetByUsername(username Name) (*User, error)

	// Find returns the set of users and total count
	// with respect to the given page request.
	Find(page PageRequest) ([]User, Count, error)

	// Add adds a new user.
	Add(user *User) (*User, error)

	// Update updates the given user.
	Update(user *User) (*User, error)

	// Delete deletes the given user.
	Delete(user *User) error
}
