package domain

// Count represents a generic count of objects.
//
// @ValueObject
type Count int64

// PageOffset represents the type for a page offset.
//
// @ValueType
type PageOffset uint

// PageSize represents the type for a page size.
//
// @ValueType
type PageSize uint

// PageRequest defines a struct for declaring pagin information for requests.
//
// @ValueObject
type PageRequest struct {
	Offset PageOffset
	Size   PageSize
}
