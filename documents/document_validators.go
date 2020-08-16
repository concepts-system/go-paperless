package documents

// import (
// 	"time"

// 	"github.com/concepts-system/go-paperless/api"
// )

// // DocumentModelValidator defines the validation rules for documents models.
// type DocumentModelValidator struct {
// 	Title         string     `json:"title" validate:"required,max=255"`
// 	Date          *time.Time `json:"date"`
// 	documentModel DocumentModel
// }

// // Bind binds the given request to a user model.
// func (v *DocumentModelValidator) Bind(c api.Context) error {
// 	if err := c.BindAndValidate(v); err != nil {
// 		return err
// 	}

// 	v.documentModel.Title = v.Title
// 	v.documentModel.Date = v.Date

// 	return nil
// }

// // NewDocumentModelValidator constructs a validator with default values.
// func NewDocumentModelValidator() DocumentModelValidator {
// 	return DocumentModelValidator{}
// }

// // NewDocumentModelValidatorFillWith constructs a validator with the values from the given document.
// func NewDocumentModelValidatorFillWith(documentModel DocumentModel) DocumentModelValidator {
// 	validator := NewDocumentModelValidator()
// 	validator.documentModel = documentModel
// 	validator.Title = documentModel.Title
// 	validator.Date = documentModel.Date

// 	return validator
// }
