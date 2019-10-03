package common

const (
	// DefaultPageSize defines the default page size.
	DefaultPageSize = 10

	// MaxPageSize defines the maximal allowed page size.
	MaxPageSize = 100
)

// PageRequest defines a struct for declaring pagin information for requests.
type PageRequest struct {
	Offset int `form:"offset"`
	Size   int `form:"size"`
}
