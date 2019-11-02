package documents

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/google/uuid"
	"github.com/kpango/glg"

	"github.com/concepts-system/go-paperless/errors"
)

// ImageConverter defines a type for bundling image conversion related operations.
type ImageConverter struct{}

const (
	mimeTypeImagePng = "image/png"
	fileExtensionPng = "png"
)

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
		return updatePageAfterConversion(page, page.Checksum)
	}

	// Convert to temp file first
	fromFile := GetContentPath(page.DocumentID, page.FileName())
	tempFileName := fmt.Sprintf("%s.%s", uuid.New(), fileExtensionPng)
	tempFile := GetContentPath(page.DocumentID, tempFileName)

	if err = i.convert(fromFile, tempFile); err != nil {
		return err
	}

	// Use temp file to obtain checksum
	file, err := os.Open(tempFile)
	if err != nil {
		return errors.Wrapf(err, "Failed to open page content file '%s'", tempFile)
	}

	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return errors.Wrapf(err, "Failed to calculate checksum on '%s'", tempFile)
	}

	checksum := fmt.Sprintf("%x", hash.Sum(nil))
	finalFileName := fmt.Sprintf("%s.%s", checksum, fileExtensionPng)
	finalFile := GetContentPath(page.DocumentID, finalFileName)

	if err = os.Rename(tempFile, finalFile); err != nil {
		return errors.Wrapf(err, "Failed to rename final page content file '%s' to '%s'", tempFile)
	}

	if err := updatePageAfterConversion(page, checksum); err != nil {
		return err
	}

	if err = os.Remove(fromFile); err != nil {
		glg.Warnf("Error during page conversion: Could not remove old page content at '%s'", fromFile)
	}

	return nil
}

func updatePageAfterConversion(page *PageModel, checksum string) error {
	page.Checksum = checksum
	page.ContentType = mimeTypeImagePng
	page.FileExtension = fileExtensionPng
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
	// TODO Make conversion command configurable
	glg.Debugf("Converting %s to %s...", from, to)
	cmd := exec.Command("convert", from, to)
	return cmd.Run()
}
