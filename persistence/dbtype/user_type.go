package dbtype

import (
	api "github.com/dankobg/juicer/api/gen"
)

type ListUsersFilters struct {
	Page     *api.PaginationPage
	PageSize *api.PaginationPageSize
	Sort     *[]string
}
