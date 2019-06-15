package documents

import (
	"github.com/jinzhu/gorm"

	"github.com/concepts-system/go-paperless/common"
	"github.com/concepts-system/go-paperless/database"
	"github.com/concepts-system/go-paperless/errors"
)

// PageState defines a string type for representing a page's state.
type PageState string

const (
	// PageStateDirty marks a page as dirty (out of sync).
	PageStateDirty = "DIRTY"
	// PageStatePending marks a page as ready for further pipeline steps.
	PageStatePending = "PENDING"
	// PageStateRecognized marks a page as recognized (OCR complete).
	PageStateRecognized = "RECOGNIZED"
	// PageStateIndexed marks a page as recognized and indexed for searching.
	PageStateIndexed = "INDEXED"
	// PageStateClean marks a page as clean (in sync).
	PageStateClean = "CLEAN"
)

// PageModel defines the data model for pages of DocumentModels.
type PageModel struct {
	gorm.Model
	DocumentID  uint      `gorm:"not_null;unique_index:idx_document_page"`
	PageNumber  uint      `gorm:"not_null;unique_index:idx_document_page"`
	State       PageState `gorm:"not_null"`
	Text        string    `gorm:"size:8192"`
	ContentType string    `gorm:"not_null;size:64"`
	ContentID   string    `gorm:"size:255"`
	Document    DocumentModel
}

// TableName for PageModel entities.
func (PageModel) TableName() string {
	return "document_pages"
}

// GetPageByDocumentIDAndPageNumber tries to find the page with the given ID and page number.
func GetPageByDocumentIDAndPageNumber(documentID uint, pageNumber uint) (*PageModel, error) {
	var document PageModel
	err := database.DB().
		Where("document_id = ? AND page_number = ?", documentID, pageNumber).
		First(&document).
		Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NotFound.Newf(
				"Page with document ID '%d' and page number '%d' not found",
				documentID,
				pageNumber,
			)
		}

		return nil, errors.Wrap(err, "Failed to fetch page")
	}

	return &document, nil
}

// GetAllPagesByDocumentID finds all pages related to the document with the given ID.
func GetAllPagesByDocumentID(documentID uint) ([]PageModel, error) {
	var pages []PageModel
	err := database.DB().
		Where("document_id = ?", documentID).
		Order("page_number").
		Find(&pages).
		Error

	if err != nil {
		return nil, err
	}

	return pages, nil
}

// FindPagesByDocumentID finds all pages related to the document with the given ID using paging.
func FindPagesByDocumentID(documentID uint, page common.PageRequest) ([]PageModel, int64, error) {
	var (
		pages      []PageModel
		totalCount int64
	)

	err := database.DB().
		Where("document_id = ?", documentID).
		Order("page_number").
		Offset(page.Offset).
		Limit(page.Size).
		Find(&pages).Error

	database.DB().Model(DocumentModel{}).Count(&totalCount)

	if err != nil {
		return nil, -1, err
	}

	return pages, totalCount, nil
}

// FindPageByID tries to find the page with the given ID.
func FindPageByID(id uint) (*PageModel, error) {
	var page PageModel
	err := database.DB().First(&page, id).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NotFound.Newf("Page with ID '%d' not found", id)
		}

		return nil, errors.Wrap(err, "Failed to fetch page")
	}

	return &page, nil
}

// Create persists the given PageModel instance in the database.
func (p *PageModel) Create() error {
	return database.DB().Create(p).Error
}

// Save saves (creates or updates) the given PageModel instance in the database.
func (p *PageModel) Save() error {
	return database.DB().Save(p).Error
}
