package domain

type DocumentSearchResult struct {
	Document *Document
}

// DocumentIndex abstracts all functionality required for indexing and searching document and pages.
type DocumentIndex interface {
	// IndexAllDocuments reinserts all documents into the index.
	IndexAllDocuments() error

	// IndexDocument inserts or updates the index entry for the document with the given document number.
	IndexDocument(documentNumber DocumentNumber) error

	// Search returns all matching documents with respect to the given query.
	Search(query string, pr PageRequest) ([]DocumentSearchResult, Count, error)
}
