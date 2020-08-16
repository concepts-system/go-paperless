package domain

import "time"

type (
	// DocumentType represents the type of a document.
	DocumentType string

	// DocumentState represents the state of a document.
	DocumentState string

	// DocumentNumber represents the of a document's unique identifier.
	DocumentNumber uint

	// ContentKey represents the type for a key pointing to a documents content.
	ContentKey string
)

const (
	// DocumentTypePDF represents the type of documents having a PDF as artifact.
	DocumentTypePDF DocumentType = "PDF"
)

const (
	// DocumentStateEdited marks a document edited.
	// A document in this state has edited pages but is not ready for
	// further processing, as not all pages have been processed.
	DocumentStateEdited = "EDITED"

	// DocumentStateDirty marks a document processed, meaning all pages are
	// processed so that the document is ready for further processing.
	DocumentStateDirty = "PROCESSED"

	// DocumentStateIndexed marks a document as indexed.
	DocumentStateIndexed = "INDEXED"

	// DocumentStateArchived marks a document as archived (in sync).
	DocumentStateArchived = "ARCHIVED"
)

// Document represents a document managed by the system.
type Document struct {
	DocumentNumber DocumentNumber
	Owner          *User
	Title          Text
	Date           *time.Time
	State          DocumentState
	Fingerprint    Fingerprint
	Type           DocumentType
	Pages          []DocumentPage
}

// AreAllPagesInState returns a boolean value indicating whether all the
// given documents pages are in the given page state.
func (d Document) AreAllPagesInState(state PageState) bool {
	for _, page := range d.Pages {
		if page.State != state {
			return false
		}
	}

	return true
}
