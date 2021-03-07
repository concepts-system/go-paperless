package web

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/concepts-system/go-paperless/domain"
)

type documentResponse struct {
	DocumentNumber uint       `json:"documentNumber,omitempty"`
	Title          string     `json:"title"`
	Date           *time.Time `json:"date,omitempty"`
	State          string     `json:"state"`
	Fingerprint    string     `json:"fingerprint,omitempty"`
	Type           string     `json:"type,omitempty"`
	PageCount      int        `json:"pageCount"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt,omitempty"`
}

type documentSearchResultResponse struct {
	Document documentResponse `json:"document"`
}

type (
	documentSerializer struct {
		C echo.Context
		*domain.Document
	}

	documentListSerializer struct {
		C         echo.Context
		Documents []domain.Document
	}

	documentSearchResultSerializer struct {
		C echo.Context
		*domain.DocumentSearchResult
	}

	documentSearchResultListSerializer struct {
		C             echo.Context
		SearchResults []domain.DocumentSearchResult
	}
)

// Response returns the API response for a document.
func (s documentSerializer) Response() documentResponse {
	return documentResponse{
		DocumentNumber: uint(s.DocumentNumber),
		Title:          string(s.Title),
		Date:           s.Date,
		Fingerprint:    string(s.Fingerprint),
		State:          string(s.State),
		Type:           string(s.Type),
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
		PageCount:      len(s.Pages),
	}
}

// Response returns the API response for a list of documents.
func (s documentListSerializer) Response() []interface{} {
	response := make([]interface{}, len(s.Documents))

	for i, document := range s.Documents {
		serializer := documentSerializer{s.C, &document}
		response[i] = serializer.Response()
	}

	return response
}

// Response returns the API response for a document search result.
func (s documentSearchResultSerializer) Response() documentSearchResultResponse {
	return documentSearchResultResponse{
		Document: documentSerializer{s.C, s.Document}.Response(),
	}
}

// Response returns the API response for a list of document search results.
func (s documentSearchResultListSerializer) Response() []interface{} {
	response := make([]interface{}, len(s.SearchResults))

	for i, result := range s.SearchResults {
		serializer := documentSearchResultSerializer{s.C, &result}
		response[i] = serializer.Response()
	}

	return response
}
