package reactive

import (
	"fmt"
	"github.com/joamaki/goreactive/stream"
)

// MapDerefPtr takes an observable of a pointer and maps it to an observable of a de-referenced value of that pointer
func MapDerefPtr[T any](pointerX stream.Observable[*T]) stream.Observable[T] {
	return stream.FlatMap(pointerX, func(ptr *T) stream.Observable[T] {
		if ptr == nil {
			return stream.Error[T](fmt.Errorf("attempted to dereference a nil ptr"))
		}
		return stream.Just(*ptr)
	})
}
