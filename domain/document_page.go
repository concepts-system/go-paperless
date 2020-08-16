package domain

type (
	// PageNumber represents the type of a document page's number.
	PageNumber uint

	// PageState represents the state of a document's page.
	PageState string

	// PageType represents the type of a document page's content.
	PageType string
)

const (
	// PageStateEdited marks a page as edited (out of sync).
	PageStateEdited = "EDITED"

	// PageStatePreprocessed marks a page as preprocessed.
	PageStatePreprocessed = "PREPROCESSED"

	// PageStateAnalyzed marks a page as analyzed (OCR complete).
	PageStateAnalyzed = "ANALYZED"

	// PageStateIndexed marks a page as recognized and indexed for searching.
	PageStateIndexed = "INDEXED"
)

// DocumentPage represents a page of a document managed by the system.
type DocumentPage struct {
	PageNumber  PageNumber
	State       PageState
	Content     Text
	Type        PageType
	Fingerprint Fingerprint
}
