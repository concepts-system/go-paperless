package documents

// import (
// 	"os"
// 	"path"
// 	"strconv"
// 	"time"

// 	"github.com/blevesearch/bleve"
// 	"github.com/blevesearch/bleve/mapping"
// 	"github.com/concepts-system/go-paperless/common"
// 	log "github.com/kpango/glg"
// )

// // Indexer represents an interface for objects able to execute index operations.
// // This creates an abstraction alyer for batches and the index itself.
// type Indexer interface {
// 	Index(id string, document interface{}) error
// }

// // DocumentIndex defines the structure of the search index for documents.
// type DocumentIndex struct {
// 	DocumentID uint
// 	OwnerID    uint
// 	Title      string
// 	Date       string
// 	Created    string
// 	Updated    string
// 	PageCount  int
// 	Pages      []PageIndex
// }

// // PageIndex defines the structure of the search index for pages.
// type PageIndex struct {
// 	PageID     uint
// 	PageNumber uint
// 	Text       string
// }

// const indexName = "document-idx"

// var index bleve.Index

// // Type returns the document type of a document index.
// func (d *DocumentIndex) Type() string {
// 	return "document"
// }

// // Type returns the document type of a page index.
// func (p *PageIndex) Type() string {
// 	return "page"
// }

// // PrepareIndex prepares the document index for usage through `Index`.
// func PrepareIndex() error {
// 	indexPath := getIndexPath()

// 	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
// 		log.Warnf("Index at '%s' not found, recreating...", indexPath)

// 		i, err := bleve.New(indexPath, createIndexMapping())
// 		if err != nil {
// 			return err
// 		}

// 		index = i
// 		log.Info("Reindexing all documents in background...")
// 		go indexAllDocuments()
// 		return nil
// 	}

// 	i, err := bleve.Open(indexPath)
// 	if err != nil {
// 		return err
// 	}

// 	index = i
// 	return nil
// }

// // GetIndex returns the singleton document index. Make sure to prepare it before using this method through `PrepareIndex`.
// func GetIndex() bleve.Index {
// 	return index
// }

// // IndexDocument indexes the document for the given document ID.
// func IndexDocument(documentID uint, indexer Indexer) error {
// 	log.Infof("Indexing document with ID %d", documentID)

// 	document, err := GetDocumentByID(documentID)
// 	if err != nil {
// 		return err
// 	}

// 	pages, err := GetAllPagesByDocumentID(documentID)
// 	if err != nil {
// 		return err
// 	}

// 	err = indexer.Index(strconv.Itoa(int(documentID)), indexForDocumentModel(document, pages))
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func createIndexMapping() mapping.IndexMapping {
// 	mapping := bleve.NewIndexMapping()
// 	// mapping.DefaultDateTimeParser = time.RFC3339
// 	mapping.AddDocumentMapping("document", createDocumentMapping())
// 	return mapping
// }

// func createDocumentMapping() *mapping.DocumentMapping {
// 	mapping := bleve.NewDocumentMapping()
// 	mapping.AddFieldMappingsAt("DocumentID", bleve.NewNumericFieldMapping())
// 	mapping.AddFieldMappingsAt("OwnerID", bleve.NewNumericFieldMapping())
// 	mapping.AddFieldMappingsAt("Title", bleve.NewTextFieldMapping())
// 	mapping.AddFieldMappingsAt("Date", bleve.NewDateTimeFieldMapping())
// 	mapping.AddFieldMappingsAt("PageCount", bleve.NewNumericFieldMapping())
// 	mapping.AddSubDocumentMapping("Pages", createPageMapping())
// 	return mapping
// }

// func createPageMapping() *mapping.DocumentMapping {
// 	mapping := bleve.NewDocumentMapping()
// 	mapping.AddFieldMappingsAt("PageID", bleve.NewNumericFieldMapping())
// 	mapping.AddFieldMappingsAt("PageNumber", bleve.NewNumericFieldMapping())
// 	mapping.AddFieldMappingsAt("Text", bleve.NewTextFieldMapping())
// 	return mapping
// }

// func getIndexPath() string {
// 	return path.Join(
// 		common.Config().GetDataPath(),
// 		indexName,
// 	)
// }

// func indexAllDocuments() {
// 	documentIDs, err := GetAllDocumentIDs()
// 	if err != nil {
// 		log.Errorf("Failed to retrieve document IDs to index: %s", err)
// 		return
// 	}

// 	batch := index.NewBatch()
// 	for _, id := range documentIDs {
// 		err := IndexDocument(id, batch)

// 		if err != nil {
// 			log.Warnf("Failed to index document with ID %d: %s", id, err)
// 		}
// 	}

// 	if err = GetIndex().Batch(batch); err != nil {
// 		log.Errorf("Failed to execute batch indexing: %s", err)
// 		return
// 	}
// }

// func indexForDocumentModel(document *DocumentModel, pages []PageModel) DocumentIndex {
// 	pageIndexes := make([]PageIndex, len(pages))
// 	for i, page := range pages {
// 		pageIndexes[i] = indexForPageModel(page)
// 	}

// 	return DocumentIndex{
// 		DocumentID: document.ID,
// 		OwnerID:    document.OwnerID,
// 		Title:      document.Title,
// 		Date:       document.Date.Format(time.RFC3339),
// 		Created:    document.CreatedAt.Format(time.RFC3339),
// 		Updated:    document.UpdatedAt.Format(time.RFC3339),
// 		PageCount:  len(pages),
// 		Pages:      pageIndexes,
// 	}
// }

// func indexForPageModel(page PageModel) PageIndex {
// 	return PageIndex{
// 		PageID:     page.ID,
// 		PageNumber: page.PageNumber,
// 		Text:       page.Text,
// 	}
// }
