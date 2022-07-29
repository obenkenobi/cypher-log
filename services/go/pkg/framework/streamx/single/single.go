package single

import (
	"context"
	"github.com/barweiss/go-tuple"
	"github.com/joamaki/goreactive/stream"
	"github.com/obenkenobi/cypher-log/services/go/pkg/framework/streamx"
)

// Single An interface that listens for a single value. It runs on top of an observable.
type Single[T any] struct {
	src stream.Observable[T]
}

func (s Single[T]) ToObservable() stream.Observable[T] {
	return s.src
}

// Just creates a single from a single item
func Just[T any](val T) Single[T] { return fromObservable(stream.Just(val)) }

// Error creates a single that fails immediately with the given error
func Error[T any](err error) Single[T] { return fromObservable(stream.Error[T](err)) }

// FromSupplier
//creates a single out of a supplier function that returns a value or an
//error to emit. If the supplier is successful, the result value is emitted.
//Otherwise, if an error is returned, an error id emitted.
func FromSupplier[T any](supplier func() (result T, err error)) Single[T] {
	var src stream.Observable[T] = stream.FuncObservable[T](func(ctx context.Context, next func(T) error) error {
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
	return fromObservable(src)
}

// FromSupplierAsync
//creates a single out of a supplier function that returns a value or an
//error to be emitted asynchronously. If the supplier is successful, the result
//value is emitted. Otherwise, if an error is returned, an error id emitted. The
//supplier function is able to emit its return value(s) asynchronously by
//running it in a separate goroutine.
func FromSupplierAsync[T any](supplier func() (result T, err error)) Single[T] {
	ch := make(chan tuple.T2[T, error])
	go func() {
		defer close(ch)
		val, err := supplier()
		ch <- tuple.New2(val, err)
	}()
	return FromSupplier[T](func() (T, error) {
		res := <-ch
		return res.V1, res.V2
	})
}

// Map applies a function onto a Single
func Map[A any, B any](src Single[A], apply func(A) B) Single[B] {
	return fromObservable[B](stream.Map(src.ToObservable(), func(a A) B { return apply(a) }))
}

func MapWithError[A any, B any](src Single[A], apply func(A) (B, error)) Single[B] {
	return fromObservable[B](
		stream.FuncObservable[B](func(ctx context.Context, next func(B) error) error {
			return src.ToObservable().Observe(
				ctx,
				func(a A) error {
					if res, err := apply(a); err != nil {
						return err
					} else {
						return next(res)
					}
				})
		}),
	)
}

func Zip[A any, B any](src1 Single[A], src2 Single[B]) Single[tuple.T2[A, B]] {
	return FlatMap(src1, func(a A) Single[tuple.T2[A, B]] {
		return Map(src2, func(b B) tuple.T2[A, B] {
			return tuple.New2(a, b)
		})
	})
}

func AwaitItem[T any](ctx context.Context, src Single[T]) (T, error) {
	return stream.First(ctx, src.ToObservable())
}

// FlatMap applies a function that returns a single of B to the source single of A.
// The Single from the function is flattened (hence FlatMap).
func FlatMap[A any, B any](src Single[A], apply func(A) Single[B]) Single[B] {
	return Single[B]{
		stream.FlatMap(src.ToObservable(), func(a A) stream.Observable[B] { return apply(a).ToObservable() }),
	}
}

// MapDerefPtr takes a single of a pointer and maps it to a single of a de-referenced value of that pointer
func MapDerefPtr[T any](src Single[*T]) Single[T] {
	return fromObservable(streamx.MapDerefPtr(src.ToObservable()))
}

func fromObservable[T any](src stream.Observable[T]) Single[T] {
	return Single[T]{src: src}
}
