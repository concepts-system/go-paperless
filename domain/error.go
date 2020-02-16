package domain

// Error represents an error occurring in the domain logic.
type Error struct {
	message string
}

func (err Error) Error() string {
	return err.message
}
