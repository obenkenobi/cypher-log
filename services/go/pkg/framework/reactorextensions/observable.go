package reactorextensions

import (
	"context"
	"fmt"
	"github.com/joamaki/goreactive/stream"
)

// ObserveSupplier creates an observable out of a supplier function that may
// return a value or an error. If the supplier is sucessful, an observable
// of that value will be created. Otherwise, if an error is returned,
// an error observable is created.
func ObserveSupplier[V any](supplier func() (V, error)) stream.Observable[V] {
	return stream.FuncObservable[V](func(ctx context.Context, next func(V) error) error {
		if ctx.Err() != nil {
			// Context already cancelled, stop before emitting items.
			return ctx.Err()
		}
		val, err := supplier()
		if err != nil {
			return err
		}
		return next(val)
	})
}

// MapDerefPtr takes an observable of a pointer and maps it to an observable of a de-referenced value of that pointer
func MapDerefPtr[V any](pointerX stream.Observable[*V]) stream.Observable[V] {
	return stream.FlatMap(pointerX, func(ptr *V) stream.Observable[V] {
		if ptr == nil {
			return stream.Error[V](fmt.Errorf("attempted to dereference a nil ptr"))
		}
		return stream.Just(*ptr)
	})
}
