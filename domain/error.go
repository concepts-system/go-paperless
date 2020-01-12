package domain

// Error represents an error occurring in the domain logic.
//
// @ValueType
type Error struct {
	message string
}

func (err Error) Error() string {
	return err.message
}
