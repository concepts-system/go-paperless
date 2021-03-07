package web

import (
	"github.com/labstack/echo/v4"

	"github.com/concepts-system/go-paperless/domain"
)

type documentPageResponse struct {
	PageNumber  uint   `json:"pageNumber,omitempty"`
	State       string `json:"state"`
	Fingerprint string `json:"fingerprint"`
	Type        string `json:"type"`
	Text        string `json:"text,omitempty"`
}

type (
	documentPageSerializer struct {
		C           echo.Context
		includeText bool
		*domain.DocumentPage
	}

	documentPageListSerializer struct {
		C     echo.Context
		Pages []domain.DocumentPage
	}
)

// Response returns the API response for a document page.
func (s documentPageSerializer) Response() documentPageResponse {
	var text string

	if s.includeText {
		text = string(s.Text)
	}

	return documentPageResponse{
		PageNumber:  uint(s.PageNumber),
		Fingerprint: string(s.Fingerprint),
		State:       string(s.State),
		Type:        string(s.Type),
		Text:        text,
	}
}

// Response returns the API response for a list of document pages.
func (s documentPageListSerializer) Response() []interface{} {
	response := make([]interface{}, len(s.Pages))

	for i, document := range s.Pages {
		response[i] = documentPageSerializer{s.C, false, &document}.Response()
	}

	return response
}
