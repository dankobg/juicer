package dbtype

import (
	api "github.com/dankobg/juicer/api/gen"
)

type ListGameResultStatusesFilters struct {
	// Page Page number (1-based)
	Page *api.PaginationPage `form:"page,omitempty" json:"page,omitempty"`

	// PageSize Number of items per page
	PageSize *api.PaginationPageSize `form:"page_size,omitempty" json:"page_size,omitempty"`

	// ID Filter game result statuses by ids
	ID *[]int64 `form:"id,omitempty" json:"id,omitempty"`

	// Name Filter game result statuses by name (partial match)
	Name *string `form:"name,omitempty" json:"name,omitempty"`

	// Sort Sort by fields (add prefix `-` for descending e.g. -created_at)
	Sort *[]string `form:"sort,omitempty" json:"sort,omitempty"`
}
