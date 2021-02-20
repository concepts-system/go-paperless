package infrastructure

import (
	"crypto/sha256"
	"encoding/hex"
	"image"
	_ "image/png"
	"io"

	"github.com/concepts-system/go-paperless/domain"
	log "github.com/sirupsen/logrus"
	"golang.org/x/image/tiff"
)

type documentPreprocessorImpl struct {
	documents       domain.Documents
	documentArchive domain.DocumentArchive
}

// NewDocumentPreprocessorImpl returns a new simple preprocessor using Go's
// standard packages.
func NewDocumentPreprocessorImpl(
	documents domain.Documents,
	documentArchive domain.DocumentArchive,
) domain.DocumentPreprocessor {
	return &documentPreprocessorImpl{
		documents,
		documentArchive,
	}
}

func (p *documentPreprocessorImpl) PreprocessPage(
	documentNumber domain.DocumentNumber,
	pageNumber domain.PageNumber,
) error {
	page, err := p.documents.GetPageByDocumentNumberAndPageNumber(
		documentNumber,
		pageNumber,
	)

	if err != nil {
		return err
	}

	if page.Type == domain.PageTypeTIFF {
		log.Debugf("Page already in correct format; skipping conversion for document %d page %d", documentNumber, pageNumber)
	} else {
		log.Debug("Converting page content to TIFF")
		if err := p.convertPageContentToTIFF(documentNumber, page); err != nil {
			return err
		}
	}

	fingerprint, err := p.computePageFingerPrint(documentNumber, page)
	if err != nil {
		return err
	}

	// Update the page's content key (hash)
	oldContentKey := page.ContentKey()
	page.Type = domain.PageTypeTIFF
	page.Fingerprint = domain.Fingerprint(fingerprint)
	err = p.documentArchive.MoveContent(documentNumber, oldContentKey, page.ContentKey())
	if err != nil {
		return err
	}

	_, err = p.documents.UpdatePage(documentNumber, page)
	if err != nil {
		log.Error("Failed to save the page after updating its content key; restoring old content key leaving the page without hash")
		if err = p.documentArchive.MoveContent(documentNumber, page.ContentKey(), oldContentKey); err != nil {
			log.Errorf(
				"Rollback failed; Page %d for document %d has now an invalid content key! This needs to be fixed manually.",
				documentNumber,
				pageNumber,
			)
			return err
		}

		log.Info("Rollback of content key successful")
		return err
	}

	log.Debug("Preprocessing done")
	return nil
}

func (p *documentPreprocessorImpl) convertPageContentToTIFF(
	documentNumber domain.DocumentNumber,
	page *domain.DocumentPage,
) error {
	content, err := p.documentArchive.ReadContent(
		documentNumber,
		page.ContentKey(),
	)

	if err != nil {
		return err
	}

	image, _, err := image.Decode(content)
	if err != nil {
		return err
	}

	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		_ = tiff.Encode(pw, image, &tiff.Options{})
	}()

	if err := p.documentArchive.StoreContent(documentNumber, page.ContentKey(), pr); err != nil {
		log.Error(err)
	}

	log.Debug("Successfully converted page content to TIFF")
	return nil
}

func (p *documentPreprocessorImpl) computePageFingerPrint(
	documentNumber domain.DocumentNumber,
	page *domain.DocumentPage,
) (string, error) {
	content, err := p.documentArchive.ReadContent(
		documentNumber,
		page.ContentKey(),
	)

	if err != nil {
		return "", err
	}

	hasher := sha256.New()
	if _, err := io.Copy(hasher, content); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
