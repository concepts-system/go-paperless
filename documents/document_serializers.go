package documents

import (
	"time"

	"github.com/labstack/echo/v4"
)

// DocumentResponse defines the document model projection returned by API methods.
type (
	DocumentResponse struct {
		ID      uint          `json:"id"`
		OwnerID uint          `json:"ownerId"`
		Title   string        `json:"title"`
		Date    *time.Time    `json:"date"`
		State   DocumentState `json:"state"`
	}

	// PageResponse defines the page model projection returned by API methods.
	PageResponse struct {
		ID         uint `json:"id"`
		DocumentID uint `json:"documentId"`
		PageNumber uint `json:"pageNumber"`
	}

	// DocumentSerializer defines functionality for serializing document models into document responses.
	DocumentSerializer struct {
		C echo.Context
		*DocumentModel
	}

	// DocumentListSerializer defines functionality for serializing document models into documents responses.
	DocumentListSerializer struct {
		C         echo.Context
		Documents []DocumentModel
	}

	// PageSerializer defines functionality for serializing page models into page responses.
	PageSerializer struct {
		C    echo.Context
		Page *PageModel
	}

	// PageListSerializer defines functionality for serializing page models into pages responses.
	PageListSerializer struct {
		C     echo.Context
		Pages []PageModel
	}
)

// Response returns the API response for a given user model.
func (s *DocumentSerializer) Response() DocumentResponse {
	return DocumentResponse{
		ID:      s.ID,
		OwnerID: s.OwnerID,
		Title:   s.Title,
		Date:    s.Date,
		State:   s.State,
	}
}

// Response returns the API response for a list of document models.
func (s *DocumentListSerializer) Response() []DocumentResponse {
	response := make([]DocumentResponse, len(s.Documents))

	for idx, document := range s.Documents {
		serializer := DocumentSerializer{s.C, &document}
		response[idx] = serializer.Response()
	}

	return response
}

// Response returns the API response for a given page model.
func (s *PageSerializer) Response() PageResponse {
	return PageResponse{
		ID:         s.Page.ID,
		DocumentID: s.Page.DocumentID,
		PageNumber: s.Page.PageNumber,
	}
}

// Response returns the API response for a list of page models.
func (s *PageListSerializer) Response() []PageResponse {
	response := make([]PageResponse, len(s.Pages))

	for idx, page := range s.Pages {
		serializer := PageSerializer{s.C, &page}
		response[idx] = serializer.Response()
	}

	return response
}
