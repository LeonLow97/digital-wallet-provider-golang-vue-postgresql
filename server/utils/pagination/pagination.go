package pagination

import (
	"math"

	"github.com/LeonLow97/go-clean-architecture/utils/constants"
)

type Paginator struct {
	TotalRecords int64 `json:"-"`
	PageSize     int64 `json:"pageSize" form:"pageSize"`
	Page         int64 `json:"page" form:"page"`
}

// Offset calculates the OFFSET value for SQL queries
func (p *Paginator) Offset() int64 {
	return (p.Page - 1) * p.PageSize
}

// Limit returns the LIMIT value for SQL queries
func (p *Paginator) Limit() int64 {
	return p.PageSize
}

// TotalPages calculates the number of pages based on total records and page size
func (p *Paginator) TotalPages() int64 {
	return int64(math.Ceil(float64(p.TotalRecords) / float64(p.PageSize)))
}

// Normalization adjusts the page size according to max and default limits
func (p *Paginator) SanitizePaginator() {
	if p.Page < 1 {
		p.Page = 1
	}

	if p.PageSize > constants.MAX_PAGE_SIZE {
		p.PageSize = constants.MAX_PAGE_SIZE
	}

	if p.PageSize < 1 {
		p.PageSize = constants.DEFAULT_PAGE_SIZE
	}
}

// HasNextPage checks if there is a next page
func (p *Paginator) HasNextPage() bool {
	return p.Page < p.TotalPages()
}

// HasPreviousPage checks it there is a previous page
func (p *Paginator) HasPreviousPage() bool {
	return p.Page > 1
}
