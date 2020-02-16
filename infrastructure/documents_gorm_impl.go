package infrastructure

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/concepts-system/go-paperless/domain"
	"github.com/concepts-system/go-paperless/errors"
)

type documentsGormImpl struct {
	db     *Database
	mapper *documentsGormMapper
}

type documentModel struct {
	DocumentNumber uint `gorm:"primary_key"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time `gorm:"index"`

	OwnerID     uint
	Title       string     `gorm:"not_null;size:255"`
	Date        *time.Time `gorm:"index"`
	State       string     `gorm:"not_null;size:32"`
	Fingerprint string     `gorm:"size:255"`
	Type        string     `gorm:"not_null;size:32"`

	Owner *userModel
	Pages []documentPageModel `gorm:"foreignkey:DocumentNumber"`
}

type documentPageModel struct {
	gorm.Model
	DocumentNumber uint   `gorm:"not_null;unique_index:idx_document_page"`
	PageNumber     uint   `gorm:"not_null;unique_index:idx_document_page"`
	State          string `gorm:"not_null;size:32"`
	Type           string `gorm:"not_null;size:32"`
	Fingerprint    string `gorm:"not_null;size:32"`
	Content        string `gorm:"size:8192"`
}

func (documentModel) TableName() string {
	return "documents"
}

func (documentPageModel) TableName() string {
	return "document_pages"
}

// NewDocuments creates a new documents domain repository.
func NewDocuments(db *Database) domain.Documents {
	return documentsGormImpl{
		db:     db,
		mapper: newDocumentsGormMapper(newUsersGormMapper()),
	}
}

func (d documentsGormImpl) FindByUsername(
	username domain.Name,
	page domain.PageRequest,
) ([]domain.Document, domain.Count, error) {
	var (
		documents  []documentModel
		totalCount int64
	)

	err := d.db.
		Joins("inner join users on users.id = documents.owner_id").
		Where("users.username = ?", username).
		Offset(page.Offset).
		Limit(page.Size).
		Find(&documents).
		Count(&totalCount).
		Error

	if err != nil {
		return nil, -1, err
	}

	return d.mapper.MapDocumentModelsToDomainEntities(documents), domain.Count(totalCount), nil
}

func (d documentsGormImpl) GetByDocumentNumber(documentNumber domain.DocumentNumber) (*domain.Document, error) {
	document, err := d.getDocumentModelByDocumentNumber(uint(documentNumber))

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return d.mapper.MapDocumentModelToDoaminEntity(document), nil
}

func (d documentsGormImpl) Add(document *domain.Document) (*domain.Document, error) {
	var owner userModel
	err := d.db.
		Select("id").
		Where("username = ?", string(document.Owner.Username)).
		First(&owner).
		Error

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to find user with username '%s'", string(document.Owner.Username))
	}

	documentModel := d.mapper.MapDomainEntityToDocumentModel(owner.ID, document)
	if err := d.db.Create(documentModel).Scan(documentModel).Error; err != nil {
		return nil, errors.Wrap(err, "Failed to create document")
	}

	return d.mapper.MapDocumentModelToDoaminEntity(documentModel), nil
}

func (d *documentsGormImpl) getDocumentModelByDocumentNumber(
	documentNumber uint,
) (*documentModel, error) {
	var document documentModel

	err := d.db.
		Preload("Owner").
		Preload("Pages").
		First(&document, documentNumber).
		Error

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &document, nil
}
