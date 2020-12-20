package infrastructure

import (
	"github.com/concepts-system/go-paperless/domain"
)

type documentPreprocessorImpl struct {
	documents       domain.Documents
	documentArchive domain.DocumentArchive
}

// NewDocumentPreprocessorImpl returns a new simple preprocessor using Go's
// standard packages.
func NewDocumentPreprocessorImpl(
	documents domain.Documents,
	documentArchive domain.DocumentArchive,
) domain.DocumentPreprocessor {
	return &documentPreprocessorImpl{
		documents,
		documentArchive,
	}
}

func (p *documentPreprocessorImpl) PreprocessPage(
	documentNumber domain.DocumentNumber,
	pageNumber domain.PageNumber,
) error {
	page, err := p.documents.GetPageByDocumentNumberAndPageNumber(
		documentNumber,
		pageNumber,
	)

	if err != nil {
		return err
	}

	_, err = p.documentArchive.ReadContent(
		documentNumber,
		page.ContentKey(),
	)

	if err != nil {
		return err
	}

	// TODO: Implement me

	return nil
}
