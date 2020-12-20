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
}

type (
	documentPageSerializer struct {
		C echo.Context
		*domain.DocumentPage
	}

	documentPageListSerializer struct {
		C     echo.Context
		Pages []domain.DocumentPage
	}
)

// Response returns the API response for a document page.
func (s documentPageSerializer) Response() documentPageResponse {
	return documentPageResponse{
		PageNumber:  uint(s.PageNumber),
		Fingerprint: string(s.Fingerprint),
		State:       string(s.State),
		Type:        string(s.Type),
	}
}

// Response returns the API response for a list of document pages.
func (s documentPageListSerializer) Response() []documentPageResponse {
	response := make([]documentPageResponse, len(s.Pages))

	for i, document := range s.Pages {
		serializer := documentPageSerializer{s.C, &document}
		response[i] = serializer.Response()
	}

	return response
}
