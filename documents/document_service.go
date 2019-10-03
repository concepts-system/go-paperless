package documents

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"github.com/concepts-system/go-paperless/common"
	"github.com/concepts-system/go-paperless/errors"
	"github.com/google/uuid"
)

const documentDataDirectoryName = "documents"

// SearchDocuments searches the document index and returns the matching documents owned by the user
// with the given user ID.
func SearchDocuments(
	userID uint,
	queryString string,
	page common.PageRequest,
	sort string,
	highlight string,
) ([]common.SearchResult, uint64, error) {
	userDocumentQuery := bleve.NewQueryStringQuery(fmt.Sprintf("OwnerID:%d", userID))

	var query query.Query
	if len(strings.TrimSpace(queryString)) > 0 {
		query = bleve.NewQueryStringQuery(queryString)
	} else {
		query = bleve.NewMatchAllQuery()
	}

	request := bleve.NewSearchRequest(bleve.NewConjunctionQuery(userDocumentQuery, query))
	request.From = page.Offset
	request.Size = page.Size

	if len(highlight) > 0 {
		request.Highlight = bleve.NewHighlightWithStyle(highlight)
	}

	if len(sort) > 0 {
		request.SortBy([]string{sort})
	}

	results, err := GetIndex().Search(request)

	if err != nil {
		return nil, 0, errors.Wrap(err, "Failed to search document index")
	}

	return common.ToSearchResults(results), results.Total, nil
}

// AppendPageToDocument appends a new page to the given document and triggers the pipeline for that document.
func AppendPageToDocument(document *DocumentModel, contentType string, content io.Reader) (*PageModel, error) {
	page, err := createPage(document.ID, contentType, content)

	if err != nil {
		return nil, err
	}

	document.State = DocumentStateDirtyPending
	if err = document.Save(); err != nil {
		return nil, err
	}

	if err = submitPageConversionJob(page.ID); err != nil {
		return nil, err
	}

	return page, nil
}

// GetContentPath returns the full path for the file containing a page's or document's content.
func GetContentPath(documentID uint, contentID string) string {
	pageContentPath := path.Join(
		common.Config().GetDataPath(),
		documentDataDirectoryName,
		fmt.Sprintf("%d", documentID),
		contentID,
	)

	os.MkdirAll(path.Dir(pageContentPath), os.ModePerm)
	return pageContentPath
}

// OpenContent opens the content file for the given document ID and content ID.
func OpenContent(documentID uint, contentID string) (*os.File, error) {
	path := GetContentPath(documentID, contentID)
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// DeleteContent removes a content file for the given document and content ID.
func DeleteContent(documentID uint, contentID string) error {
	path := GetContentPath(documentID, contentID)
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return os.Remove(path)
	}

	return nil
}

// Helper methods

func createPage(documentID uint, contentType string, content io.Reader) (*PageModel, error) {
	extension, err := common.GetExtensionByMimeType(contentType)
	if err != nil {
		return nil, err
	}

	pages, err := GetAllPagesByDocumentID(documentID)
	if err != nil {
		return nil, err
	}

	contentID := fmt.Sprintf("%s%s", uuid.New(), extension)
	page := PageModel{
		DocumentID:  documentID,
		PageNumber:  uint(len(pages)),
		State:       PageStateDirty,
		ContentType: contentType,
		ContentID:   contentID,
	}

	contentPath := GetContentPath(documentID, contentID)
	file, err := os.Create(contentPath)

	if err != nil {
		return nil, err
	}

	defer file.Close()
	io.Copy(file, content)

	if err = page.Create(); err != nil {
		return nil, err
	}

	return &page, nil
}
