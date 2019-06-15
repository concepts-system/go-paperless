package documents

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/google/uuid"
	"github.com/kpango/glg"
)

// ImageConverter defines a type for bundling image conversion related operations.
type ImageConverter struct{}

const mimeTypeImagePng = "image/png"

var (
	typesToConvert = []string{
		"image/bmp",
		"image/tiff",
	}
)

// ConvertPage convertss the page with the given ID to PNG for all formats that are marked for conversion.
func (i ImageConverter) ConvertPage(pageID uint) error {
	page, err := FindPageByID(pageID)
	if err != nil {
		return errors.New("Page not found")
	}

	if !i.needsConversion(page.ContentType) {
		glg.Infof("Skipping conversion for page %d as content type %s does not require conversion", pageID, page.ContentType)
		return updatePageAfterConversion(page, page.ContentID)
	}

	fromFile := GetContentPath(page.DocumentID, page.ContentID)
	newContentID := fmt.Sprintf("%s.png", uuid.New())
	toFile := GetContentPath(page.DocumentID, newContentID)

	if err = i.convert(fromFile, toFile); err != nil {
		return err
	}

	if err := updatePageAfterConversion(page, newContentID); err != nil {
		return err
	}

	if err = os.Remove(fromFile); err != nil {
		glg.Warnf("Error during page conversion: Could not remove old page content at %s", fromFile)
	}

	return nil
}

func updatePageAfterConversion(page *PageModel, contentID string) error {
	page.ContentID = contentID
	page.ContentType = mimeTypeImagePng
	page.State = PageStatePending

	return page.Save()
}

func (i ImageConverter) needsConversion(contentType string) bool {
	for _, typ := range typesToConvert {
		if typ == contentType {
			return true
		}
	}

	return false
}

func (i ImageConverter) convert(from, to string) error {
	// TODO: Make conversion command configurable
	glg.Debugf("Converting %s to %s...", from, to)
	cmd := exec.Command("convert", from, to)
	return cmd.Run()
}
