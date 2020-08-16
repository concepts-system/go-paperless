package documents

// import (
// 	"fmt"
// 	"time"

// 	"github.com/jinzhu/gorm"

// 	"github.com/concepts-system/go-paperless/common"
// 	"github.com/concepts-system/go-paperless/database"
// 	"github.com/concepts-system/go-paperless/errors"
// 	"github.com/concepts-system/go-paperless/users"
// )

// // DocumentState defines a string type for representating a document's state.
// type DocumentState string

// const (
// 	// DocumentStateDirtyPending marks a document as dirty/pending.
// 	// A document in this state has dirty pages but is not ready for generation,
// 	// as not all pages are fully indexed.
// 	DocumentStateDirtyPending = "EDITED"

// 	// DocumentStateDirty marks a document as dirty (out of sync and ready for generation).
// 	DocumentStateDirty = "PROCESSED"

// 	// DocumentStateIndexed marks a document as indexed.
// 	DocumentStateIndexed = "INDEXED"

// 	// DocumentStateClean marks a document as archived (in sync).
// 	DocumentStateClean = "ARCHIVED"
// )

// // DocumentModel defines the data model for documents managed by the system.
// type DocumentModel struct {
// 	gorm.Model
// 	OwnerID       uint            `gorm:"not_null"`
// 	Owner         users.UserModel `gorm:"not_null"`
// 	Title         string          `gorm:"not_null;size:255"`
// 	Date          *time.Time      `gorm:"index"`
// 	State         DocumentState   `gorm:"not_null;size:32"`
// 	ContentID     string          `gorm:"size:255"`
// 	FileExtension string          `gorm:"size:8"`
// 	Pages         []PageModel
// }

// // TableName for DocumentModel entities.
// func (DocumentModel) TableName() string {
// 	return "documents"
// }

// // FileName returns the name this page's data is stored under on the file system.
// func (d DocumentModel) FileName() string {
// 	return fmt.Sprintf("%s.%s", d.ContentID, d.FileExtension)
// }

// // GetDocumentByID tries to find the document with the given ID.
// func GetDocumentByID(id uint) (*DocumentModel, error) {
// 	var document DocumentModel
// 	err := database.DB().First(&document, id).Error

// 	if err != nil {
// 		if gorm.IsRecordNotFoundError(err) {
// 			return nil, errors.NotFound.Newf("Document with ID '%d' not found", id)
// 		}

// 		return nil, errors.Wrap(err, "Failed to fetch document")
// 	}

// 	return &document, nil
// }

// // GetAllDocumentIDs returns an array of all document IDs from the database.
// func GetAllDocumentIDs() ([]uint, error) {
// 	var (
// 		documents []DocumentModel
// 		ids       []uint
// 	)

// 	err := database.DB().Find(&documents).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	ids = make([]uint, len(documents))
// 	for i, document := range documents {
// 		ids[i] = document.ID
// 	}

// 	return ids, nil
// }

// // FindDocumentsByUserID returns documents within the system based on the given offset and limit params.
// // The method will return the total count of found objects as 2nd parameter.
// func FindDocumentsByUserID(userID uint, page common.PageRequest) ([]DocumentModel, int64, error) {
// 	var (
// 		documents  []DocumentModel
// 		totalCount int64
// 	)

// 	err := database.DB().
// 		Where("owner_id = ?", userID).
// 		Order("date desc").
// 		Offset(page.Offset).
// 		Limit(page.Size).
// 		Find(&documents).Error

// 	database.DB().Model(DocumentModel{}).Count(&totalCount)

// 	if err != nil {
// 		return nil, -1, err
// 	}

// 	return documents, totalCount, nil
// }

// // Create persists the given DocumentModel instance in the database.
// func (d *DocumentModel) Create() error {
// 	return database.DB().Create(d).Error
// }

// // Save saves (creates or updates) the given DocumentModel instance in the database.
// func (d *DocumentModel) Save() error {
// 	return database.DB().Save(d).Error
// }
