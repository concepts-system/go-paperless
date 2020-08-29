package domain

// DocumentPreprocessor defines the signature of a component being capable
// of preprocessing a document page to preper for further operations.
type DocumentPreprocessor interface {
	// PreprocessPage applies preprocessing to a document's page.
	PreprocessPage(documentNumber DocumentNumber, pageNumber PageNumber) error
}
