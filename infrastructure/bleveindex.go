package infrastructure

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/concepts-system/go-paperless/common"
	"github.com/concepts-system/go-paperless/domain"
	"github.com/concepts-system/go-paperless/errors"
	"github.com/sirupsen/logrus"
)

const indexBatchSize = 100

type indexer interface {
	Index(id string, document interface{}) error
}

type bleveIndex struct {
	logger    *logrus.Entry
	documents domain.Documents
	index     bleve.Index
}

type documentEntry struct {
	DocumentNumber uint
	OwnerUsername  string
	Title          string
	Date           *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	PageCount      int
	Pages          []pageEntry
}

type pageEntry struct {
	PageNumber uint
	Text       string
}

func (d *documentEntry) Type() string {
	return "document"
}

func (p *pageEntry) Type() string {
	return "document_page"
}

// NewBleveDocumentIndex returns a document index implementation based on a local Bleve index.
func NewBleveDocumentIndex(
	documentIndexPath string,
	documents domain.Documents,
) (domain.DocumentIndex, error) {
	documentIndex := &bleveIndex{
		logger:    common.NewLogger("bleve-index"),
		documents: documents,
	}

	if err := documentIndex.initializeDocumentIndex(documentIndexPath); err != nil {
		return nil, err
	}

	return documentIndex, nil
}

func (b *bleveIndex) IndexAllDocuments() error {
	b.logger.Info("Indexing all documents")
	pr := domain.PageRequest{Size: indexBatchSize}

	for {
		documents, totalCount, err := b.documents.Find(pr)
		if err != nil {
			return err
		}

		if len(documents) == 0 {
			break
		}

		batch := b.index.NewBatch()
		for _, document := range documents {
			err := b.indexDocument(document, batch)

			if err != nil {
				return errors.Wrapf(err, "Failed to index document %d", document.DocumentNumber)
			}
		}

		if err := b.index.Batch(batch); err != nil {
			return errors.Wrapf(err, "Failed to index document batch")
		}

		b.logger.Debugf("Indexed documents: %d / %d", pr.Offset+len(documents), totalCount)
		pr.Offset += indexBatchSize
	}

	b.logger.Info("Indexing complete")
	return nil
}

func (b *bleveIndex) IndexDocument(documentNumber domain.DocumentNumber) error {
	document, err := b.documents.GetByDocumentNumber(documentNumber)
	if err != nil {
		return err
	}

	return b.indexDocument(*document, b.index)
}

func (b *bleveIndex) Search(
	queryString string,
	page domain.PageRequest,
) ([]domain.DocumentSearchResult, domain.Count, error) {
	var query query.Query
	if len(strings.TrimSpace(queryString)) > 0 {
		query = bleve.NewQueryStringQuery(queryString)
	} else {
		query = bleve.NewMatchAllQuery()
	}

	request := bleve.NewSearchRequest(query)
	request.From = page.Offset
	request.Size = page.Size
	if len(page.Sort) > 0 {
		request.SortBy([]string{page.Sort})
	}

	results, err := b.index.Search(request)
	if err != nil {
		return nil, 0, errors.Wrap(err, "Failed to search document index")
	}

	return b.mapDocumentSearchResults(results), domain.Count(results.Total), nil
}

/* Helper Methods */

func (b *bleveIndex) initializeDocumentIndex(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		b.logger.Warnf("Index at '%s' not found, creating new index", path)
		b.index, err = bleve.New(path, b.createIndexMapping())
		if err != nil {
			return err
		}

		go func() {
			err := b.IndexAllDocuments()
			if err != nil {
				b.logger.Error(err)
			}
		}()

		return nil
	}

	index, err := bleve.Open(path)
	if err != nil {
		return err
	}

	b.index = index
	return nil
}

func (b *bleveIndex) createIndexMapping() mapping.IndexMapping {
	mapping := bleve.NewIndexMapping()
	// mapping.DefaultDateTimeParser = time.RFC3339

	mapping.AddDocumentMapping("document", b.createDocumentIndexMapping())

	return mapping
}

func (b *bleveIndex) createDocumentIndexMapping() *mapping.DocumentMapping {
	mapping := bleve.NewDocumentMapping()

	mapping.AddFieldMappingsAt("DocumentNumber", bleve.NewNumericFieldMapping())
	mapping.AddFieldMappingsAt("OwnerUsername", bleve.NewTextFieldMapping())
	mapping.AddFieldMappingsAt("Title", bleve.NewTextFieldMapping())
	mapping.AddFieldMappingsAt("Date", bleve.NewDateTimeFieldMapping())
	mapping.AddFieldMappingsAt("UpdatedAt", bleve.NewDateTimeFieldMapping())
	mapping.AddFieldMappingsAt("CreatedAt", bleve.NewDateTimeFieldMapping())
	mapping.AddFieldMappingsAt("PageCount", bleve.NewNumericFieldMapping())
	mapping.AddSubDocumentMapping("Pages", b.createDocumentPageIndexMapping())

	return mapping
}

func (b *bleveIndex) createDocumentPageIndexMapping() *mapping.DocumentMapping {
	mapping := bleve.NewDocumentMapping()

	mapping.AddFieldMappingsAt("PageNumber", bleve.NewNumericFieldMapping())
	mapping.AddFieldMappingsAt("Text", bleve.NewTextFieldMapping())

	return mapping
}

func (b *bleveIndex) indexDocument(document domain.Document, indexer indexer) error {
	return indexer.Index(fmt.Sprint(document.DocumentNumber), b.documentEntry(document))
}

func (b *bleveIndex) documentEntry(document domain.Document) *documentEntry {
	entry := documentEntry{
		DocumentNumber: uint(document.DocumentNumber),
		OwnerUsername:  string(document.Owner.Username),
		Title:          string(document.Title),
		Date:           document.Date,
		CreatedAt:      document.CreatedAt,
		UpdatedAt:      document.UpdatedAt,
		PageCount:      len(document.Pages),
		Pages:          make([]pageEntry, len(document.Pages)),
	}

	for i, page := range document.Pages {
		entry.Pages[i] = *b.documentPageEntry(page)
	}

	return &entry
}

func (b *bleveIndex) documentPageEntry(page domain.DocumentPage) *pageEntry {
	return &pageEntry{
		PageNumber: uint(page.PageNumber),
		Text:       string(page.Text),
	}
}

func (b *bleveIndex) mapDocumentSearchResults(result *bleve.SearchResult) []domain.DocumentSearchResult {
	results := make([]domain.DocumentSearchResult, len(result.Hits))

	for i, hit := range result.Hits {
		documentNumber, _ := strconv.Atoi(hit.ID)
		results[i] = domain.DocumentSearchResult{
			DocumentNumber: domain.DocumentNumber(documentNumber),
		}
	}

	return results
}
