package web

import (
	"time"

	"github.com/concepts-system/go-paperless/domain"
)

type documentValidator struct {
	Title string     `json:"title" validate:"required,min=1,max=255"`
	Date  *time.Time `json:"date"`

	document domain.Document
}

// Bind binds the given request to a document model.
func (v *documentValidator) Bind(c *context) error {
	if err := c.BindAndValidate(v); err != nil {
		return err
	}

	v.document.Title = domain.Text(v.Title)
	v.document.Date = v.Date

	return nil
}

func newDocumentValidator() *documentValidator {
	return &documentValidator{}
}

// func newDocumentValidatorOf(document *domain.Document) *documentValidator {
// 	validator := newDocumentValidator()
// 	validator.Title = string(document.Title)
// 	validator.Date = document.Date

// 	return validator
// }
