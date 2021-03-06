package domain

import (
	"fmt"
	"strings"
)

type (
	// PageNumber represents the type of a document page's number.
	PageNumber uint

	// PageState represents the state of a document's page.
	PageState string

	// PageType represents the type of a document page's content.
	PageType string
)

const (
	// PageTypeTIFF
	PageTypeTIFF = PageType("TIFF")
	// PageTypeUnknown
	PageTypeUnknown = PageType("UNKNOWN")
)

const (
	// PageStateEdited marks a page as edited (out of sync).
	PageStateEdited = PageState("EDITED")

	// PageStatePreprocessed marks a page as preprocessed.
	PageStatePreprocessed = PageState("PREPROCESSED")

	// PageStateAnalyzed marks a page as analyzed (OCR complete).
	PageStateAnalyzed = PageState("ANALYZED")

	// PageStateIndexed marks a page as recognized and indexed for searching.
	PageStateIndexed = PageState("INDEXED")
)

// DocumentPage represents a page of a document managed by the system.
type DocumentPage struct {
	PageNumber  PageNumber
	State       PageState
	Text        Text
	Type        PageType
	Fingerprint Fingerprint
	IsInReview  bool
	Document    *Document
}

// ContentKey returns the content key for the document.
func (d DocumentPage) ContentKey() ContentKey {
	return ContentKey(fmt.Sprintf(
		"%s.%s",
		d.Fingerprint,
		strings.ToLower(string(d.Type)),
	))
}
