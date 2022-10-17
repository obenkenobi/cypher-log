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

type PageRequest struct {
	Page int
	Size int
	Sort []SortField
}

func NewPageRequest(page int, size int) PageRequest {
	return PageRequest{Page: page, Size: size}
}

type Page[T any] struct {
	Contents []T
	Total    int
}

func NewPage[T any](contents []T, total int) Page[T] {
	return Page[T]{Contents: contents, Total: total}
}

func Map[T any, V any](page Page[T], handler func(T) V) Page[V] {
	return NewPage(slice.Map(page.Contents, handler), page.Total)
}
