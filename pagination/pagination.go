package pagination

import api "github.com/dankobg/juicer/api/gen"

var (
	DefaultPage     = 1
	DefaultPageSize = 50
	MinPage         = 1
	MinPageSize     = 1
	MaxPageSize     = 500
)

type CursorPagination struct {
	PageSize int
	Cursor   string
}

type Pagination struct {
	Page     int
	PageSize int
}

type PaginationMeta struct {
	Pagination
	Total      int64
	HasMore    bool
	NextCursor string
}

func NewMeta(np Pagination, totalCount int64) PaginationMeta {
	return PaginationMeta{
		Pagination: Pagination{
			Page:     np.Page,
			PageSize: np.PageSize,
		},
		Total: totalCount,
	}
}

func (pm PaginationMeta) ToResp() api.PaginationMeta {
	return api.PaginationMeta{
		Page:     pm.Page,
		PageSize: pm.PageSize,
		Total:    pm.Total,
	}
}

func GetNormalized(page *int, pageSize *int) Pagination {
	p := Pagination{Page: DefaultPage, PageSize: DefaultPageSize}
	if page != nil {
		p.Page = max(*page, MinPage)
	}

	if pageSize != nil {
		p.PageSize = min(max(*pageSize, MinPageSize), MaxPageSize)
	}

	return p
}

type Result[T any] struct {
	Data []T
	Meta PaginationMeta
}

func NewRes[T any](data []T, meta PaginationMeta) Result[T] {
	return Result[T]{
		Data: data,
		Meta: meta,
	}
}

type WithTotal[T any] struct {
	Data       []T
	TotalCount int64
}

func NewWithTotal[T any](data []T, totalCount int64) WithTotal[T] {
	return WithTotal[T]{
		Data:       data,
		TotalCount: totalCount,
	}
}

type WithHasMore[T any] struct {
	Data    []T
	HasMore bool
}

func NewWithHasMore[T any](data []T, hasMore bool) WithHasMore[T] {
	return WithHasMore[T]{
		Data:    data,
		HasMore: hasMore,
	}
}
