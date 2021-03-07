package web

import (
	"github.com/concepts-system/go-paperless/domain"
)

const (
	// DefaultPageSize defines the default page size.
	defaultPageSize = 10

	// MaxPageSize defines the maximal allowed page size.
	maxPageSize = 100
)

// PageRequest defines a struct for declaring pagin information for requests.
type pageRequest struct {
	Offset int    `form:"offset" query:"offset"`
	Size   int    `form:"size" query:"offset"`
	Sort   string `form:"sort" query:"sort"`
}

// ToDomainPageRequest maps the given page request to a domain value.
func (pr pageRequest) ToDomainPageRequest() domain.PageRequest {
	return domain.PageRequest{
		Offset: domain.PageOffset(pr.Offset),
		Size:   domain.PageSize(pr.Size),
	}
}
