package domain

// Documents defines an interface for managing the collection of all documents.
type Documents interface {
	// FindByUsername returns the set of documents owned by the user
	// with the given username and total count with respect to the given
	// page request.
	FindByUsername(username Name, pr PageRequest) ([]Document, Count, error)

	// GetByDocumentNumber returns the document with the given document number
	// or nil in case no such document exists.
	GetByDocumentNumber(documentNumber DocumentNumber) (*Document, error)

	// Add adds the given document.
	Add(document *Document) (*Document, error)
}
