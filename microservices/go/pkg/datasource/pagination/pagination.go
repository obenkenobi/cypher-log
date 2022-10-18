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
	Field     string    `json:"field"`
	Direction Direction `json:"direction"`
}

func NewSortField(field string, direction Direction) SortField {
	return SortField{Field: field, Direction: direction}
}

type PageRequest struct {
	Page int64       `json:"page"`
	Size int64       `json:"size"`
	Sort []SortField `json:"sort"`
}

func (p PageRequest) SkipCount() int64 {
	return p.Page * p.Size
}

func NewPageRequest(page int64, size int64) PageRequest {
	return PageRequest{Page: page, Size: size}
}

type Page[T any] struct {
	Contents []T   `json:"contents"`
	Total    int64 `json:"total"`
}

func NewPage[T any](contents []T, total int64) Page[T] {
	return Page[T]{Contents: contents, Total: total}
}

func Map[T any, V any](page Page[T], handler func(T) V) Page[V] {
	return NewPage(slice.Map(page.Contents, handler), page.Total)
}
