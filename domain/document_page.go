package domain

type (
	// PageNumber represents the type of a document page's number.
	PageNumber uint

	// PageState represents the state of a document's page.
	PageState string

	// PageType represents the type of a document page's content.
	PageType string
)

// DocumentPage represents a page of a document managed by the system.
type DocumentPage struct {
	PageNumber  PageNumber
	State       PageState
	Content     Text
	Type        PageType
	Fingerprint Fingerprint
}
