package infrastructure

import (
	"time"

	"gorm.io/gorm"

	"github.com/concepts-system/go-paperless/domain"
	"github.com/concepts-system/go-paperless/errors"
)

type documentsGormImpl struct {
	db     *Database
	mapper *documentsGormMapper
}

type documentModel struct {
	DocumentNumber uint `gorm:"not_null;primary_key"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time `gorm:"index"`

	OwnerID     uint
	Title       string     `gorm:"not_null;size:255"`
	Date        *time.Time `gorm:"index"`
	State       string     `gorm:"not_null;size:32"`
	Fingerprint string     `gorm:"size:255;index"`
	Type        string     `gorm:"not_null;size:32"`

	Owner *userModel
	Pages []documentPageModel `gorm:"foreignkey:DocumentNumber"`
}

type documentPageModel struct {
	DocumentNumber uint `gorm:"not_null;primaryKey;autoIncrement:false"`
	PageNumber     uint `gorm:"not_null;primaryKey;autoIncrement:false"`

	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`
	State       string     `gorm:"not_null;size:32"`
	Type        string     `gorm:"not_null;size:32"`
	Fingerprint string     `gorm:"not_null;size:32"`
	Text        string     `gorm:"size:8192"`
	IsInReview  bool

	Document *documentModel `gorm:"foreignKey:DocumentNumber"`
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
		if gorm.ErrRecordNotFound == err {
			return nil, nil
		}

		return nil, err
	}

	return d.mapper.MapDocumentModelToDoaminEntity(document), nil
}

func (d documentsGormImpl) Add(document *domain.Document) (*domain.Document, error) {
	owner, err := d.getDocumentOwner(document)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create document")
	}

	documentModel := d.mapper.MapDomainEntityToDocumentModel(owner.ID, document)
	if err := d.db.Create(documentModel).Scan(documentModel).Error; err != nil {
		return nil, errors.Wrap(err, "Failed to create document")
	}

	return d.mapper.MapDocumentModelToDoaminEntity(documentModel), nil
}

func (d documentsGormImpl) Update(document *domain.Document) (*domain.Document, error) {
	_, err := d.getDocumentModelByDocumentNumber(uint(document.DocumentNumber))
	if err != nil {
		return nil, err
	}

	owner, err := d.getDocumentOwner(document)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to update document")
	}

	documentModel := d.mapper.MapDomainEntityToDocumentModel(owner.ID, document)
	if err := d.db.Save(documentModel).Scan(documentModel).Error; err != nil {
		return nil, errors.Wrap(err, "Failed to update document")
	}

	return d.mapper.MapDocumentModelToDoaminEntity(documentModel), nil
}

func (d documentsGormImpl) GetPagesByDocumentNumber(
	documentNumber domain.DocumentNumber,
	page domain.PageRequest,
) ([]domain.DocumentPage, domain.Count, error) {
	var totalCount int64
	var pageModels []documentPageModel
	err := d.db.
		Where("document_number = ?", documentNumber).
		Offset(page.Offset).
		Limit(page.Size).
		Find(&pageModels).
		Count(&totalCount).
		Error

	if err != nil {
		return nil, -1, errors.Wrapf(err, "Failed to retrieve document pages")
	}

	pages := d.mapper.MapPageModelsToDomainEntities(pageModels)
	return pages, domain.Count(totalCount), nil
}

func (d documentsGormImpl) GetPageByDocumentNumberAndPageNumber(
	documentNumber domain.DocumentNumber,
	pageNumber domain.PageNumber,
) (*domain.DocumentPage, error) {
	page, err := d.getDocumentPageModelByDocumentNumberAndPageNumber(
		uint(documentNumber),
		uint(pageNumber),
	)

	if err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, nil
		}

		return nil, err
	}

	return d.mapper.MapPageModelToDomainEntity(page), nil
}

func (d documentsGormImpl) AddPage(
	documentNumber domain.DocumentNumber,
	page *domain.DocumentPage,
) (*domain.DocumentPage, error) {
	document, err := d.getDocumentModelByDocumentNumber(uint(documentNumber))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create page")
	}

	pageModel := d.mapper.MapDomainEntityToPageModel(document.DocumentNumber, page)
	if err := d.db.Create(pageModel).Scan(pageModel).Error; err != nil {
		return nil, errors.Wrap(err, "Failed to create page")
	}

	return d.mapper.MapPageModelToDomainEntity(pageModel), nil
}

func (d documentsGormImpl) UpdatePage(
	documentNumber domain.DocumentNumber,
	page *domain.DocumentPage,
) (*domain.DocumentPage, error) {
	document, err := d.getDocumentModelByDocumentNumber(uint(documentNumber))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to update page")
	}

	pageModel := d.mapper.MapDomainEntityToPageModel(document.DocumentNumber, page)
	if err := d.db.Save(pageModel).Scan(pageModel).Error; err != nil {
		return nil, errors.Wrap(err, "Failed to update page")
	}

	return d.mapper.MapPageModelToDomainEntity(pageModel), nil
}

/* Helper Methods */

func (d *documentsGormImpl) getDocumentOwner(document *domain.Document) (*userModel, error) {
	var owner userModel
	err := d.db.
		Select("id").
		Where("username = ?", string(document.Owner.Username)).
		First(&owner).
		Error

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to find user with username '%s'", string(document.Owner.Username))
	}

	return &owner, nil
}

func (d *documentsGormImpl) getDocumentModelByDocumentNumber(
	documentNumber uint,
) (*documentModel, error) {
	document := documentModel{
		DocumentNumber: documentNumber,
	}

	err := d.db.
		Preload("Owner").
		Preload("Pages").
		First(&document).
		Error

	if err != nil {
		return nil, err
	}

	return &document, nil
}

func (d *documentsGormImpl) getDocumentPageModelByDocumentNumberAndPageNumber(
	documentNumber uint,
	pageNumber uint,
) (*documentPageModel, error) {
	documentPage := documentPageModel{
		DocumentNumber: documentNumber,
		PageNumber:     pageNumber,
	}

	err := d.db.Preload("Document").First(&documentPage).Error

	if err != nil {
		return nil, err
	}

	return &documentPage, nil
}
