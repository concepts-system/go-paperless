package domain

import "time"

type (
	// DocumentType represents the type of a document.
	DocumentType string

	// DocumentState represents the state of a document.
	DocumentState string

	// DocumentNumber represents the of a document's unique identifier.
	DocumentNumber uint
)

const (
	// DocumentTypePDF represents the type of documents having a PDF as artifact.
	DocumentTypePDF DocumentType = "PDF"
)

const (
	// DocumentStateNew represents the state new documents are in.
	DocumentStateNew DocumentState = "NEW"
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
