package infrastructure

import (
	"bytes"
	"io"
	"os/exec"

	"github.com/concepts-system/go-paperless/domain"
	"github.com/concepts-system/go-paperless/errors"
	log "github.com/sirupsen/logrus"
)

const (
	tesseractExecutable = "tesseract"
	languages           = "eng+deu"
)

// TesseractOcrEngine provides an interface to the Tesseract OCR engine.
type TesseractOcrEngine struct {
	documents       domain.Documents
	documentArchive domain.DocumentArchive
}

// NewTesseractOcrEngine returns a new Tesseract OCR engine.
func NewTesseractOcrEngine(
	documents domain.Documents,
	documentArchive domain.DocumentArchive,
) *TesseractOcrEngine {
	return &TesseractOcrEngine{
		documents,
		documentArchive,
	}
}

// RecognizePage executes OCR to get the text from the pages image content.
func (t *TesseractOcrEngine) ScanPage(
	documentNumber domain.DocumentNumber,
	pageNumber domain.PageNumber,
) error {
	page, err := t.documents.GetPageByDocumentNumberAndPageNumber(documentNumber, pageNumber)
	if err != nil {
		return err
	}

	content, err := t.documentArchive.ReadContent(page.Document.DocumentNumber, page.ContentKey())
	if err != nil {
		return err
	}

	result, err := t.recognizeImage(content)
	if err != nil {
		return errors.Wrapf(
			err,
			"Scanning failed for document '%d' page '%d'",
			page.Document.DocumentNumber,
			page.PageNumber,
		)
	}

	page.Text = domain.Text((result.Bytes()))
	page.State = domain.PageStateAnalyzed

	_, err = t.documents.UpdatePage(documentNumber, page)
	content.Close()
	return err
}

/* Helper Functions */

func (t *TesseractOcrEngine) recognizeImage(reader io.Reader) (*bytes.Buffer, error) {
	path, err := exec.LookPath(tesseractExecutable)
	if err != nil {
		return nil, err
	}

	buffer := &bytes.Buffer{}
	cmd := exec.Cmd{
		Path:   path,
		Args:   []string{tesseractExecutable, "-l", languages, "stdin", "stdout"},
		Stdin:  reader,
		Stdout: buffer,
		Stderr: log.StandardLogger().Out,
	}

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return buffer, nil
}
