package postgres

import (
	"reflect"
	"slices"
	"strings"

	api "github.com/dankobg/juicer/api/gen"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
	"github.com/stephenafamo/bob/dialect/psql/sm"
)

func AddPagination(q *bob.BaseQuery[*dialect.SelectQuery], page *api.PaginationPage, pageSize *api.PaginationPageSize) {
	if page == nil && pageSize == nil {
		return
	}

	q.Apply(sm.Columns(psql.Raw("COUNT(*) OVER()").As("total_count")))

	var limit, offset int64
	if pageSize != nil && *pageSize > 0 {
		limit = int64(*pageSize)
		q.Apply(sm.Limit(limit))
	}

	if page != nil && *page > 0 {
		if limit > 0 {
			offset = (int64(*page) - 1) * limit
		} else {
			offset = int64(*page)
		}

		q.Apply(sm.Offset(offset))
	}
}

func AddOrderBy(q *bob.BaseQuery[*dialect.SelectQuery], sort *[]string, allowedCols []string) {
	if sort == nil {
		return
	}

	for _, field := range *sort {
		field = strings.TrimSpace(field)
		direction := "ASC"
		column := field

		if strings.HasPrefix(field, "-") {
			direction = "DESC"
			column = strings.TrimPrefix(field, "-")
		}

		if !slices.Contains(allowedCols, column) {
			continue
		}

		if direction == "ASC" {
			q.Apply(sm.OrderBy(column).Asc())
		} else {
			q.Apply(sm.OrderBy(column).Desc())
		}
	}
}

// HasAnyFilters returns true if at least one filter was provided (not nil), unless excluded
func HasAnyFilters[T any](params *T, excludeFields ...string) bool {
	if params == nil {
		return false
	}

	excludeMap := make(map[string]struct{}, len(excludeFields))
	for _, f := range excludeFields {
		excludeMap[f] = struct{}{}
	}

	v := reflect.ValueOf(params).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := v.Type().Field(i).Name

		if _, excluded := excludeMap[fieldName]; excluded {
			continue
		}

		if field.Kind() == reflect.Pointer && !field.IsNil() {
			return true
		}
	}

	return false
}

// HasAnyLogicFilters returns true if any filter that isn't pagination or sort is provided
func HasAnyLogicFilters[T any](params *T) bool {
	return HasAnyFilters(params, "Page", "PageSize", "Sort")
}
