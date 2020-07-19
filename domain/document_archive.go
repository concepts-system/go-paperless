package domain

import (
	"io"
)

// DocumentArchive provides an abstraction for document storages.
// They provide functionality for storing and retreiving document page data.
type DocumentArchive interface {
	// ReadContent returns the content of a document or page from the store.
	ReadContent(documentNumber DocumentNumber, contentKey ContentKey) (io.Reader, error)

	// StoreContent stores content for a new document or page.
	StoreContent(documentNumber DocumentNumber, contentKey ContentKey, content io.Reader) error

	// MoveContent moves the stored content of one key to another.
	MoveContent(documentNumber DocumentNumber, sourceContentKey ContentKey, destinationContentKey ContentKey) error

	// DeleteContent deletes the content for a document or page from the store.
	DeleteContent(documentNumber DocumentNumber, contentKey ContentKey) error
}
