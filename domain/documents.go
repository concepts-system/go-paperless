package domain

// Documents defines an interface for managing the collection of all documents.
type Documents interface {
	// Find returns the a subset of all documents with respect to the given page request.
	Find(pr PageRequest) ([]Document, Count, error)

	// FindByUsername returns the set of documents owned by the user
	// with the given username, alongside with the total count with respect to the given
	// page request.
	FindByUsername(username Name, pr PageRequest) ([]Document, Count, error)

	// GetByDocumentNumber returns the document with the given document number
	// or nil in case no such document exists.
	GetByDocumentNumber(documentNumber DocumentNumber) (*Document, error)

	// Add adds the given document without its pages.
	Add(document *Document) (*Document, error)

	// Update updates the given document without its pages.
	Update(document *Document) (*Document, error)

	// GetPagesByDocumentNumber returns all pages contained in the document for the given document number
	// alongside the total count of pages with respect to the given page request.
	GetPagesByDocumentNumber(documentNumber DocumentNumber, pr PageRequest) ([]DocumentPage, Count, error)

	// GetDocumentPageByDocumentNumberAndPageNumber returns the given document
	// page with the given page number, part of the document with the given
	// document number.
	GetPageByDocumentNumberAndPageNumber(
		documentNumber DocumentNumber,
		pageNumber PageNumber,
	) (*DocumentPage, error)

	// AddPage adds the given page to the document with the given document
	// number.
	AddPage(
		document DocumentNumber,
		page *DocumentPage,
	) (*DocumentPage, error)

	// UpdatePage saves the given document page associated to the given
	// document.
	UpdatePage(
		documentNumber DocumentNumber,
		page *DocumentPage,
	) (*DocumentPage, error)
}
