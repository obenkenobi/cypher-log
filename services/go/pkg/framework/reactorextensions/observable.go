package reactorextensions

import (
	"context"
	"fmt"
	"github.com/barweiss/go-tuple"
	"github.com/joamaki/goreactive/stream"
)

// ObserveSupplier
//creates an observable out of a supplier function that returns a value or an
//error to emit. If the supplier is successful, the result value is emitted.
//Otherwise, if an error is returned, an error id emitted.
func ObserveSupplier[V any](supplier func() (result V, err error)) stream.Observable[V] {
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

// ObserveSupplierAsync
//creates an observable out of a supplier function that returns a value or an
//error to be emitted asynchronously. If the supplier is successful, the result
//value is emitted. Otherwise, if an error is returned, an error id emitted. The
//supplier function is able to emit its return value(s) asynchronously by
//running it in a separate goroutine.
func ObserveSupplierAsync[V any](supplier func() (result V, err error)) stream.Observable[V] {
	ch := make(chan tuple.T2[V, error])
	go func() {
		defer close(ch)
		val, err := supplier()
		ch <- tuple.New2(val, err)
	}()
	return ObserveSupplier[V](func() (V, error) {
		res := <-ch
		return res.V1, res.V2
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
