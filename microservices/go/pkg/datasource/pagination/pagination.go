package pagination

import (
	"github.com/akrennmair/slice"
)

type Direction string

const (
	Ascending  Direction = "asc"
	Descending Direction = "desc"
)

type SortField struct {
	Field     string
	Direction Direction
}

func NewSortField(field string, direction Direction) SortField {
	return SortField{Field: field, Direction: direction}
}

type PageRequest struct {
	Page int64
	Size int64
	Sort []SortField
}

func (p PageRequest) SkipCount() int64 {
	return p.Page * p.Size
}

func NewPageRequest(page int64, size int64) PageRequest {
	return PageRequest{Page: page, Size: size}
}

type Page[T any] struct {
	Contents []T
	Total    int64
}

func NewPage[T any](contents []T, total int64) Page[T] {
	return Page[T]{Contents: contents, Total: total}
}

func Map[T any, V any](page Page[T], handler func(T) V) Page[V] {
	return NewPage(slice.Map(page.Contents, handler), page.Total)
}
