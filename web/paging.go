package web

import (
	"github.com/concepts-system/go-paperless/domain"
)

const (
	// DefaultPageSize defines the default page size.
	DefaultPageSize = 10

	// MaxPageSize defines the maximal allowed page size.
	MaxPageSize = 100
)

// PageRequest defines a struct for declaring pagin information for requests.
type PageRequest struct {
	Offset int `form:"offset"`
	Size   int `form:"size"`
}

// ToDomainPageRequest maps the given page request to a domain value.
func (pr PageRequest) ToDomainPageRequest() domain.PageRequest {
	return domain.PageRequest{
		Offset: domain.PageOffset(pr.Offset),
		Size:   domain.PageSize(pr.Size),
	}
}
