package domain

import (
	"errors"

	"github.com/labstack/gommon/log"
)

const (
	mailboxPagePreprocess = Mailbox("document.page.preprocess")
	malboxPageAnalyze     = Mailbox("document.page.analyze")
)

// DocumentRegistry provides an abstraction for document components taking
// care of the overall document workflow.
type DocumentRegistry interface {
	Review(documentNumber DocumentNumber)
}

type documentRegistryImpl struct {
	tubeMail     TubeMail
	documents    Documents
	preprocessor DocumentPreprocessor
}

func NewDocumentRegistry(
	tubeMail TubeMail,
	documents Documents,
	preprocessor DocumentPreprocessor,
) DocumentRegistry {
	registry := &documentRegistryImpl{
		tubeMail,
		documents,
		preprocessor,
	}

	registry.setupTubeMail()
	return *registry
}

func (d documentRegistryImpl) Review(documentNumber DocumentNumber) {
	//
}

func (d documentRegistryImpl) setupTubeMail() {
	// Document-specific receivers

	// Page-specific receivers
	d.registerDocumentPageReceiver(mailboxPagePreprocess, d.preprocessPage)
	d.registerDocumentPageReceiver(malboxPageAnalyze, d.analyzePage)
}

/* Handlers */

func (d documentRegistryImpl) preprocessPage(
	documentNumber DocumentNumber,
	pageNumber PageNumber,
) error {
	_, err := d.documents.GetPageByDocumentNumberAndPageNumber(documentNumber, pageNumber)
	if err != nil {
		return err
	}

	if err = d.preprocessor.PreprocessPage(documentNumber, pageNumber); err != nil {
		return err
	}

	page, err := d.documents.GetPageByDocumentNumberAndPageNumber(documentNumber, pageNumber)
	if err != nil {
		return err
	}

	page.State = PageStatePreprocessed
	if _, err = d.documents.UpdatePage(documentNumber, page); err != nil {
		return err
	}

	d.Review(documentNumber)
	return nil
}

func (d documentRegistryImpl) analyzePage(
	documentNumber DocumentNumber,
	pageNumber PageNumber,
) error {

}

/* Helper Methods */

func (d documentRegistryImpl) registerDocumentReceiver(
	mailbox Mailbox,
	handler func(DocumentNumber) error,
) {
	receiver := func(message ...interface{}) error {
		if len(message) != 2 {
			return errors.New("Unexpected document message length")
		}

		documentNumber, ok := message[0].(DocumentNumber)
		if !ok {
			return errors.New("Unexpected document message format")
		}

		log.Infof(
			"Received message in '%s': document %s",
			mailbox,
			documentNumber,
		)

		return handler(documentNumber)
	}

	d.tubeMail.RegisterReceiver(mailbox, receiver)
}

func (d documentRegistryImpl) registerDocumentPageReceiver(
	mailbox Mailbox,
	handler func(DocumentNumber, PageNumber) error,
) {
	receiver := func(message ...interface{}) error {
		if len(message) != 2 {
			return errors.New("Unexpected document page message length")
		}

		documentNumber, documentNumberOk := message[0].(DocumentNumber)
		pageNumber, pageNumberOk := message[1].(PageNumber)
		if !documentNumberOk || !pageNumberOk {
			return errors.New("Unexpected document page message format")
		}

		log.Infof(
			"Received message in '%s': document %s, page %s",
			mailbox,
			documentNumber,
			pageNumber,
		)

		return handler(documentNumber, pageNumber)
	}

	d.tubeMail.RegisterReceiver(mailbox, receiver)
}
