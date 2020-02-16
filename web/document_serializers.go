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
	Type           string     `json:"type"`
	PageCount      int        `json:"pageCount"`
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
		PageCount:      len(s.Pages),
	}
}

// Response returns the API response for a list of documents.
func (s documentListSerializer) Response() []documentResponse {
	response := make([]documentResponse, len(s.Documents))

	for i, document := range s.Documents {
		serializer := documentSerializer{s.C, &document}
		response[i] = serializer.Response()
	}

	return response
}
