package pagination

type Paginator struct {
	TotalRecords int64 `json:"totalRecords"`
	PageSize     int
}
