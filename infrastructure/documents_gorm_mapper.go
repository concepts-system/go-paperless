package infrastructure

import "github.com/concepts-system/go-paperless/domain"

type documentsGormMapper struct {
	usersMapper *usersGormMapper
}

func newDocumentsGormMapper(usersMapper *usersGormMapper) *documentsGormMapper {
	return &documentsGormMapper{
		usersMapper: usersMapper,
	}
}

// MapDocumentModelToDoaminEntity maps the given document model to a
// corresponding domain entity.
func (m *documentsGormMapper) MapDocumentModelToDoaminEntity(
	document *documentModel,
) *domain.Document {
	if document == nil {
		return nil
	}

	return &domain.Document{
		DocumentNumber: domain.DocumentNumber(document.DocumentNumber),
		Title:          domain.Text(document.Title),
		Date:           document.Date,
		State:          domain.DocumentState(document.State),
		Fingerprint:    domain.Fingerprint(document.Fingerprint),
		Type:           domain.DocumentType(document.Type),
		Owner:          m.usersMapper.MapUserModelToDomainEntity(document.Owner),
		Pages:          m.MapPageModelsToDomainEntities(document.Pages),
	}
}

// MapDocumentModelsToDomainEntities maps the given list of document models to a list
// containing the corresponding domain entities.
func (m *documentsGormMapper) MapDocumentModelsToDomainEntities(documents []documentModel) []domain.Document {
	if documents == nil {
		return nil
	}

	domainEntities := make([]domain.Document, len(documents))

	for i, document := range documents {
		domainEntities[i] = *m.MapDocumentModelToDoaminEntity(&document)
	}

	return domainEntities
}

// MapDomainEntityToDocumentModel maps the given domain entity to the corresponding document model.
func (m *documentsGormMapper) MapDomainEntityToDocumentModel(ownerID uint, document *domain.Document) *documentModel {
	if document == nil {
		return nil
	}

	model := &documentModel{
		DocumentNumber: uint(document.DocumentNumber),
		OwnerID:        ownerID,
		Title:          string(document.Title),
		Date:           document.Date,
		State:          string(document.State),
		Fingerprint:    string(document.Fingerprint),
		Type:           string(document.Type),
	}

	return model
}

// MapPageModelToDomainEntity maps the given page model to the corresponding domain entity.
func (m *documentsGormMapper) MapPageModelToDomainEntity(page *documentPageModel) *domain.DocumentPage {
	if page == nil {
		return nil
	}

	return &domain.DocumentPage{
		PageNumber:  domain.PageNumber(page.PageNumber),
		State:       domain.PageState(page.State),
		Content:     domain.Text(page.Content),
		Type:        domain.PageType(page.Type),
		Fingerprint: domain.Fingerprint(page.Fingerprint),
	}
}

// MapPageModelsToDomainEntities maps the given list of page models to a list
// containing the corresponding domain entities.
func (m *documentsGormMapper) MapPageModelsToDomainEntities(pages []documentPageModel) []domain.DocumentPage {
	if pages == nil {
		return nil
	}

	domainEntities := make([]domain.DocumentPage, len(pages))

	for i, page := range pages {
		domainEntities[i] = *m.MapPageModelToDomainEntity(&page)
	}

	return domainEntities
}
