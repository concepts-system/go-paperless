package domain

import (
	"errors"

	"github.com/concepts-system/go-paperless/common"
)

const (
	mailboxPagePreprocess = Mailbox("document.page.preprocess")
	malboxPageAnalyze     = Mailbox("document.page.analyze")
)

var log = common.NewLogger("registry")

// DocumentRegistry provides an abstraction for document components taking
// care of the overall document workflow.
type DocumentRegistry interface {
	Review(documentNumber DocumentNumber)
}

// Receiver defines the signature of an abstract tube mail receiver.
type Receiver = func(...interface{}) error

type documentRegistryImpl struct {
	tubeMail     TubeMail
	documents    Documents
	preprocessor DocumentPreprocessor
	analyzer     DocumentAnalyzer
}

func NewDocumentRegistry(
	tubeMail TubeMail,
	documents Documents,
	preprocessor DocumentPreprocessor,
	analyzer DocumentAnalyzer,
) DocumentRegistry {
	registry := &documentRegistryImpl{
		tubeMail,
		documents,
		preprocessor,
		analyzer,
	}

	registry.setupTubeMail()
	return *registry
}

func (d documentRegistryImpl) Review(documentNumber DocumentNumber) {
	log.Debugf("Reviewing document %d", documentNumber)
	// TODO: Make sure documents are not reviewed multiple times in parallel
	//		 i.e. add some document.isInReview flag
	document, err := d.documents.GetByDocumentNumber(documentNumber)
	if err != nil {
		log.Error(err)
	}

	switch document.State {
	case DocumentStateEmpty:
		log.Debug("Document is empty; nothing to do")
	case DocumentStateEdited:
		log.Debug("Document has been edited since last review; reviewing pages")
		d.reviewDocumentPages(document)
	case DocumentStateArchived:
		log.Debug("Document is already archived; nothing to do")
	default:
		log.Warnf("Documents in state %s are not handled yet!", document.State)
	}
}

func (d documentRegistryImpl) reviewDocumentPages(document *Document) {
	for _, page := range document.Pages {
		d.reviewDocumentPage(document.DocumentNumber, page)
	}
}

func (d documentRegistryImpl) reviewDocumentPage(documentNumber DocumentNumber, page DocumentPage) {
	var err error

	switch page.State {
	case PageStateEdited:
		log.Debug("Page has been modified; sending to preprocessing")
		err = d.tubeMail.SendMessage(mailboxPagePreprocess, documentNumber, page.PageNumber)
	case PageStatePreprocessed:
		log.Debug("Page is preprocessed; sending to scanning")
		err = d.tubeMail.SendMessage(malboxPageAnalyze, documentNumber, page.PageNumber)
	default:
		log.Warnf("Document pages in state %s are not handled yet!", page.State)
	}

	if err != nil {
		log.Error(err)
	}
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
	if err := d.preprocessor.PreprocessPage(documentNumber, pageNumber); err != nil {
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
	if err := d.analyzer.ScanPage(documentNumber, pageNumber); err != nil {
		return err
	}

	page, err := d.documents.GetPageByDocumentNumberAndPageNumber(documentNumber, pageNumber)
	if err != nil {
		return err
	}

	page.State = PageStateAnalyzed
	if _, err = d.documents.UpdatePage(documentNumber, page); err != nil {
		return err
	}

	d.Review(documentNumber)
	return nil
}

/* Helper Methods */

// func (d documentRegistryImpl) registerDocumentReceiver(
// 	mailbox Mailbox,
// 	handler func(DocumentNumber) error,
// ) {
// 	receiver := func(message ...interface{}) error {
// 		if len(message) != 2 {
// 			return errors.New("Unexpected document message length")
// 		}

// 		documentNumber, ok := message[0].(DocumentNumber)
// 		if !ok {
// 			return errors.New("Unexpected document message format")
// 		}

// 		log.Infof(
// 			"Received message in '%v': document %v",
// 			mailbox,
// 			documentNumber,
// 		)

// 		return handler(documentNumber)
// 	}

// 	d.mustRegisterReceiver(mailbox, receiver)
// }

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
			"{document: %v, page: %v} -> @%v",
			documentNumber,
			pageNumber,
			mailbox,
		)

		return handler(documentNumber, pageNumber)
	}

	d.mustRegisterReceiver(mailbox, receiver)
}

func (d documentRegistryImpl) mustRegisterReceiver(mailbox Mailbox, receiver Receiver) {
	if err := d.tubeMail.RegisterReceiver(mailbox, receiver); err != nil {
		log.Fatal("Failed to register receiver", err)
	}
}
