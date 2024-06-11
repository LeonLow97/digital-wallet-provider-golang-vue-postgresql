package pagination

import (
	"math"

	"github.com/LeonLow97/go-clean-architecture/utils/constants"
)

type Paginator struct {
	TotalRecords int64 `json:"totalRecords"`
	PageSize     int   `json:"pageSize"`
	Page         int   `json:"page"`
}

// NewPaginator creates a new Paginator instance
func NewPaginator(pageSize, page int) *Paginator {
	paginator := &Paginator{
		PageSize: pageSize,
		Page:     page,
	}

	paginator.Normalization()

	return paginator
}

// Offset calculates the OFFSET value for SQL queries
func (p *Paginator) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Limit returns the LIMIT value for SQL queries
func (p *Paginator) Limit() int {
	return p.PageSize
}

// TotalPages calculates the number of pages based on total records and page size
func (p *Paginator) TotalPages() int {
	return int(math.Ceil(float64(p.TotalRecords) / float64(p.PageSize)))
}

// Normalization adjusts the page size according to max and default limits
func (p *Paginator) Normalization() {
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
