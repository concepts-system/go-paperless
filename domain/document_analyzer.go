package domain

// DocumentAnalyzer defines functionality for obtaining text from pages.
type DocumentAnalyzer interface {
	ScanPage(documentNumber DocumentNumber, pageNumber PageNumber) error
}
