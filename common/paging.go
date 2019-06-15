package common

const (
	// DefaultPageSize defines the default page size.
	DefaultPageSize = 10
)

// PageRequest defines a struct for declaring pagin information for requests.
type PageRequest struct {
	Offset int `form:"offset"`
	Size   int `form:"size"`
}
