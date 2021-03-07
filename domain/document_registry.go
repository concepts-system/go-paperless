package domain

import (
	"errors"

	"github.com/concepts-system/go-paperless/common"
)

const (
	mailboxDocumentIndex = Mailbox("document.index")

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
	index        DocumentIndex
	analyzer     DocumentAnalyzer
}

func NewDocumentRegistry(
	tubeMail TubeMail,
	documents Documents,
	preprocessor DocumentPreprocessor,
	analyzer DocumentAnalyzer,
	index DocumentIndex,
) DocumentRegistry {
	registry := &documentRegistryImpl{
		tubeMail,
		documents,
		preprocessor,
		index,
		analyzer,
	}

	registry.setupTubeMail()
	return *registry
}

func (d documentRegistryImpl) Review(documentNumber DocumentNumber) {
	log.Debugf("Reviewing document %d", documentNumber)
	document, err := d.documents.GetByDocumentNumber(documentNumber)
	if err != nil {
		log.Error(err)
	}

	switch document.State {
	case DocumentStateEmpty:
		log.Debug("Document is empty; indexing")
		err = d.indexDocument(documentNumber)
	case DocumentStateEdited:
		log.Debug("Document has been edited since last review; reviewing pages")
		d.reviewDocumentPages(document)
	case DocumentStateArchived:
		log.Debug("Document is already archived; nothing to do")
	default:
		log.Warnf("Documents in state %s are not handled yet!", document.State)
	}

	if err != nil {
		log.Error(err)
	}
}

func (d documentRegistryImpl) reviewDocumentPages(document *Document) {
	var err error

	if document.AreAllPagesInState(PageStateAnalyzed) {
		err = d.tubeMail.SendMessage(mailboxDocumentIndex, document.DocumentNumber)
	} else {
		for _, page := range document.Pages {
			d.reviewDocumentPage(document.DocumentNumber, &page)
		}
	}

	if err != nil {
		log.Error(err)
	}
}

func (d documentRegistryImpl) reviewDocumentPage(documentNumber DocumentNumber, page *DocumentPage) {
	pageNumber := page.PageNumber

	page, err := d.startPageReview(documentNumber, page)
	if err != nil {
		log.Error(err)
	}

	if page == nil {
		log.Infof("Document %d page %d is already in review; skipping", documentNumber, pageNumber)
		return
	}

	switch page.State {
	case PageStateEdited:
		log.Debug("Page has been modified; sending to preprocessing")
		err = d.tubeMail.SendMessage(mailboxPagePreprocess, documentNumber, page.PageNumber)
	case PageStatePreprocessed:
		log.Debug("Page is preprocessed; sending to scanning")
		err = d.tubeMail.SendMessage(malboxPageAnalyze, documentNumber, page.PageNumber)
	default:
		log.Warnf("Document pages in state %s are not handled yet!", page.State)
		_, err = d.finishPageReview(documentNumber, pageNumber, page.State)
	}

	if err != nil {
		log.Error(err)
	}
}

func (d documentRegistryImpl) setupTubeMail() {
	// Document-specific receivers
	d.registerDocumentReceiver(mailboxDocumentIndex, d.indexDocument)

	// Page-specific receivers
	d.registerDocumentPageReceiver(mailboxPagePreprocess, d.preprocessPage)
	d.registerDocumentPageReceiver(malboxPageAnalyze, d.analyzePage)
}

/* Handlers */

func (d documentRegistryImpl) indexDocument(documentNumber DocumentNumber) error {
	if err := d.index.IndexDocument(documentNumber); err != nil {
		return err
	}

	if _, err := d.finishDocumentReview(documentNumber, DocumentStateIndexed); err != nil {
		return err
	}

	d.Review(documentNumber)
	return nil
}

func (d documentRegistryImpl) preprocessPage(
	documentNumber DocumentNumber,
	pageNumber PageNumber,
) error {
	if err := d.preprocessor.PreprocessPage(documentNumber, pageNumber); err != nil {
		return err
	}

	if _, err := d.finishPageReview(documentNumber, pageNumber, PageStatePreprocessed); err != nil {
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

	if _, err := d.finishPageReview(documentNumber, pageNumber, PageStateAnalyzed); err != nil {
		return err
	}

	d.Review(documentNumber)
	return nil
}

/* Helper Methods */

func (d documentRegistryImpl) registerDocumentReceiver(
	mailbox Mailbox,
	handler func(DocumentNumber) error,
) {
	receiver := func(message ...interface{}) error {
		if len(message) != 1 {
			return errors.New("Unexpected document message length")
		}

		documentNumber, ok := message[0].(DocumentNumber)
		if !ok {
			return errors.New("Unexpected document message format")
		}

		log.Infof(
			"{document: %v} -> @%v",
			documentNumber,
			mailbox,
		)

		return handler(documentNumber)
	}

	d.mustRegisterReceiver(mailbox, receiver)
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

func (d documentRegistryImpl) finishDocumentReview(documentNumber DocumentNumber, state DocumentState) (*Document, error) {
	document, err := d.documents.GetByDocumentNumber(documentNumber)
	if err != nil {
		log.Error(err)
	}

	// TODO: Introduce 'isInReview' flag as for pages
	document.State = state
	document, err = d.documents.Update(document)
	if err != nil {
		return nil, err
	}

	return document, nil
}

func (d documentRegistryImpl) startPageReview(documentNumber DocumentNumber, page *DocumentPage) (*DocumentPage, error) {
	if page.IsInReview {
		return nil, nil
	}

	page.IsInReview = true
	page, err := d.documents.UpdatePage(documentNumber, page)

	if err != nil {
		return nil, err
	}

	return page, nil
}

func (d documentRegistryImpl) finishPageReview(documentNumber DocumentNumber, pageNumber PageNumber, state PageState) (*DocumentPage, error) {
	page, err := d.documents.GetPageByDocumentNumberAndPageNumber(documentNumber, pageNumber)
	if err != nil {
		log.Error(err)
	}

	if !page.IsInReview {
		return nil, nil
	}

	page.IsInReview = false
	page.State = state
	page, err = d.documents.UpdatePage(documentNumber, page)
	if err != nil {
		return nil, err
	}

	return page, nil
}
