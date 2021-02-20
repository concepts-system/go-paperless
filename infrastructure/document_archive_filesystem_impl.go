package infrastructure

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/concepts-system/go-paperless/domain"
	"github.com/concepts-system/go-paperless/errors"
)

type documentArchiveFileSystemImpl struct {
	basePath string
}

// NewDocumentArchiveFileSystemImpl returns a new document store using file
// system as a backend for storing documents.
func NewDocumentArchiveFileSystemImpl(basePath string) (domain.DocumentArchive, error) {
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		if err = os.MkdirAll(basePath, os.ModePerm); err != nil {
			return nil, errors.Wrap(err, "Failed to create data directory")
		}
	}

	return &documentArchiveFileSystemImpl{basePath}, nil
}

// GetPageContent returns the content of a document or page from the store.
func (store *documentArchiveFileSystemImpl) ReadContent(
	documentNumber domain.DocumentNumber,
	contentKey domain.ContentKey,
) (io.Reader, error) {
	path := store.getContentPath(documentNumber, contentKey)
	if _, err := os.Stat(path); err != nil {
		return nil, errors.Wrapf(err, "Content file '%s' does not exist", path)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to open content file '%s'", path)
	}

	return file, nil
}

// StoreContent stores content for a new document or page.
func (store *documentArchiveFileSystemImpl) StoreContent(
	documentNumber domain.DocumentNumber,
	contentKey domain.ContentKey,
	content io.Reader,
) error {
	path := store.getContentPath(documentNumber, contentKey)
	directory := filepath.Dir(path)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if err := os.MkdirAll(directory, os.ModePerm); err != nil {
			return errors.Wrapf(err, "Failed to create directory '%s'", directory)
		}
	}

	file, err := os.Create(path)

	if err != nil {
		return errors.Wrapf(err, "Failed to create content file '%s'", path)
	}

	defer file.Close()
	if _, err = io.Copy(file, content); err != nil {
		return errors.Wrapf(err, "Failed to write content file '%s'", path)
	}

	return nil
}

// StoreContent stores content for a new document or page.
func (store *documentArchiveFileSystemImpl) MoveContent(
	documentNumber domain.DocumentNumber,
	sourceContentKey domain.ContentKey,
	destinationContentKey domain.ContentKey,
) error {
	sourcePath := store.getContentPath(documentNumber, sourceContentKey)
	if _, err := os.Stat(sourcePath); err != nil {
		return errors.Wrapf(err, "Source content file '%s' does not exist", sourcePath)
	}

	destinationPath := store.getContentPath(documentNumber, destinationContentKey)
	if _, err := os.Stat(destinationPath); err != nil && !os.IsNotExist(err) {
		return errors.Wrapf(err, "Destination content file '%s' does already exist", destinationPath)
	}

	return os.Rename(sourcePath, destinationPath)
}

// DeleteContent deletes the content for a document or page from the store.
func (store *documentArchiveFileSystemImpl) DeleteContent(
	documentNumber domain.DocumentNumber,
	contentKey domain.ContentKey,
) error {
	path := store.getContentPath(documentNumber, contentKey)
	if err := os.Remove(path); err != nil {
		return errors.Wrapf(err, "Failed to delete content file '%s'", path)
	}

	return nil
}

/* Helper Methods */

func (store *documentArchiveFileSystemImpl) getContentPath(
	documentNumber domain.DocumentNumber,
	contentKey domain.ContentKey,
) string {
	return path.Join(
		store.basePath,
		fmt.Sprintf("%d", uint(documentNumber)),
		string(contentKey),
	)
}
