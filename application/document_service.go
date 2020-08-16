package application

import (
	"github.com/concepts-system/go-paperless/domain"
	"github.com/concepts-system/go-paperless/errors"
)

// DocumentService defines an application service for managing document-related
// use cases.
type DocumentService interface {
	// GetUserDocuments returns the given user's documents with respect to the
	// given page request.
	GetUserDocuments(username string, pr domain.PageRequest) ([]domain.Document, int64, error)

	// GetUserDocumentByDocumentNumber returns the document with the given document number owned by the given user.
	GetUserDocumentByDocumentNumber(username string, documentNumber uint) (*domain.Document, error)

	// CreateNewDocument creates the given new document owned by the user with the given username.
	CreateNewDocument(username string, document *domain.Document) (*domain.Document, error)
}

type documentServiceImpl struct {
	users           domain.Users
	documents       domain.Documents
	documentArchive domain.DocumentArchive
}

// NewDocumentService creates a new document service.
func NewDocumentService(
	users domain.Users,
	documents domain.Documents,
	documentArchive domain.DocumentArchive,
) DocumentService {
	return &documentServiceImpl{
		users:           users,
		documents:       documents,
		documentArchive: documentArchive,
	}
}

func (s *documentServiceImpl) GetUserDocuments(
	username string,
	pr domain.PageRequest,
) ([]domain.Document, int64, error) {
	documents, count, err := s.documents.FindByUsername(domain.Name(username), pr)

	if err != nil {
		return nil, -1, errors.Wrap(err, "Failed to retreive documents")
	}

	return documents, int64(count), nil
}

func (s *documentServiceImpl) GetUserDocumentByDocumentNumber(
	username string,
	documentNumber uint,
) (*domain.Document, error) {
	document, err := s.expectDocumentWithDocumentNumberExists(domain.DocumentNumber(documentNumber))

	if err != nil {
		return nil, err
	}

	if err := s.expectUserMayAccessDocument(domain.Name(username), document); err != nil {
		return nil, err
	}

	return document, nil
}

func (s *documentServiceImpl) CreateNewDocument(username string, document *domain.Document) (*domain.Document, error) {
	owner, err := s.users.GetByUsername(domain.Name(username))
	if err != nil {
		return nil, err
	}

	document.Owner = owner
	document.State = domain.DocumentStateEdited
	document.Type = ""
	document.Fingerprint = ""
	document.Pages = nil

	newDocument, err := s.documents.Add(document)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create document")
	}

	return newDocument, nil
}

/* Helper Methods */

func (s *documentServiceImpl) userMayAccessDocument(
	username string,
	document *domain.Document,
) (bool, error) {
	return document.Owner.Username == domain.Name(username), nil
}

func (s *documentServiceImpl) expectDocumentWithDocumentNumberExists(
	documentNumber domain.DocumentNumber,
) (*domain.Document, error) {
	document, err := s.documents.GetByDocumentNumber(documentNumber)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to retrieve document")
	}

	if document == nil {
		return nil, NotFoundError.Newf("Document %d does not exist", documentNumber)
	}

	return document, nil
}

func (s *documentServiceImpl) expectUserMayAccessDocument(username domain.Name, document *domain.Document) error {
	mayAccess, err := s.userMayAccessDocument(string(username), document)
	if err != nil {
		return err
	}

	if !mayAccess {
		return ForbiddenError.Newf("Access to document %d not permitted", document.DocumentNumber)
	}

	return nil
}
