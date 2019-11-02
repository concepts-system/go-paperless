package documents

import (
	"bytes"
	"io"
	"os/exec"
	"strings"

	"github.com/concepts-system/go-paperless/errors"
	"github.com/google/uuid"
)

const fileExtensionPdf = "pdf"

// OcrEngine defines a type for bundling OCR related operations.
type OcrEngine struct{}

// RecognizePage executes OCR to get the text from the pages image content.
func (o OcrEngine) RecognizePage(page *PageModel) (*bytes.Buffer, error) {
	content, err := o.recognizeFile(GetContentPath(page.DocumentID, page.FileName()), "deu+eng")
	if err != nil {
		return nil, errors.Wrapf(err, "Recognition failed for page %d", page.ID)
	}

	return content, nil
}

// GenerateDocument generates a searchable PDF for the given model and its pages.
// Returns the new content ID and the file extension of the created PDF; may be used as document's new content ID.
func (o OcrEngine) GenerateDocument(document *DocumentModel) (string, string, error) {
	pages, err := GetAllPagesByDocumentID(document.ID)
	if err != nil {
		return "", "", err
	}

	if !allPagesClean(pages) {
		return "", "", errors.BadRequest.Newf("Generation failed since not all pages of document %d are clean", document.ID)
	}

	pagePaths := make([]string, len(pages))
	for idx, page := range pages {
		pagePaths[idx] = GetContentPath(page.DocumentID, page.FileName())
	}

	newContentID := uuid.New().String()
	if err := generateSearchablePDF(GetContentPath(document.ID, newContentID), "eng+deu", pagePaths); err != nil {
		return "", "", errors.Wrapf(err, "Failed to generate document with ID %d", document.ID)
	}

	return newContentID, fileExtensionPdf, nil
}

func (o OcrEngine) recognizeFile(path, languages string) (*bytes.Buffer, error) {
	// TODO Make OCR command/path configurable
	cmd := exec.Command("tesseract", "-l", languages, path, "-")
	var buffer bytes.Buffer
	cmd.Stdout = &buffer

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return &buffer, nil
}

func generateSearchablePDF(outputPath, languages string, pagePaths []string) error {
	inputCommand := exec.Command("echo", strings.Join(pagePaths, "\n"))
	generationCommand := exec.Command("tesseract", "-l", languages, "stdin", outputPath, "pdf")

	// inputCommand | w => r | generationCommand
	r, w := io.Pipe()
	generationCommand.Stdin, inputCommand.Stdout = r, w

	if err := inputCommand.Start(); err != nil {
		return err
	}
	if err := generationCommand.Start(); err != nil {
		return err
	}
	if err := inputCommand.Wait(); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	if err := generationCommand.Wait(); err != nil {
		return err
	}

	return nil
}
