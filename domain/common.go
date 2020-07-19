package domain

type (
	// Text represents the type of abstract text.
	Text string

	// Count represents a generic count of objects.
	Count int64

	// Fingerprint represents the type of a document's or page's fingerprint.
	Fingerprint string

	// PageOffset represents the type for a page offset.
	PageOffset uint

	// PageSize represents the type for a page size.
	PageSize uint
)

// PageRequest defines a struct for declaring pagin information for requests.
type PageRequest struct {
	Offset PageOffset
	Size   PageSize
}
