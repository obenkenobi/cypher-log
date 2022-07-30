package single

import (
	"context"
	"github.com/barweiss/go-tuple"
	"github.com/joamaki/goreactive/stream"
	"github.com/obenkenobi/cypher-log/services/go/pkg/extensions/streamx"
	"sync"
)

// Single An interface that listens for a single value. It runs on top of an observable.
type Single[T any] struct {
	src stream.Observable[T]
}

func (s Single[T]) ToObservable() stream.Observable[T] {
	return s.src
}

// Thread safe Observable that reads from a channel, used for ensuring when a
// single is made from a channel, it will always emit the same value even when
// observed multiple times. This is done by caching the channel result in a
// thread safe manner.
type channelObservable[T any] struct {
	channelRead bool
	ch          <-chan T
	value       T
	valueRWLock sync.RWMutex
}

func (a *channelObservable[T]) Observe(ctx context.Context, next func(T) error) error {
	if ctx.Err() != nil {
		// Context already cancelled, stop before emitting items.
		return ctx.Err()
	}
	a.valueRWLock.RLock()
	shouldAttemptWriteValue := !a.channelRead
	a.valueRWLock.RUnlock()

	if shouldAttemptWriteValue {
		a.valueRWLock.Lock()
		if !a.channelRead {
			a.value = <-a.ch
			a.channelRead = true
		}
		a.valueRWLock.Unlock()
	}
	a.valueRWLock.RLock()
	value := a.value
	a.valueRWLock.RUnlock()
	return next(value)
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

func FromChannel[T any](ch <-chan T) Single[T] {
	return Single[T]{
		src: &channelObservable[T]{channelRead: false, ch: ch},
	}
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
	return FlatMap(FromChannel(ch), func(res tuple.T2[T, error]) Single[T] {
		val, err := res.V1, res.V2
		if err != nil {
			return Error[T](err)
		}
		return Just(val)
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

// Zip2 Takes 2 Singles and returns a Single that emits a tuple of each of the
// singles in the order they are supplied
func Zip2[V1 any, V2 any](src1 Single[V1], src2 Single[V2]) Single[tuple.T2[V1, V2]] {
	return FlatMap(src1, func(a V1) Single[tuple.T2[V1, V2]] {
		return Map(src2, func(b V2) tuple.T2[V1, V2] { return tuple.New2(a, b) })
	})
}

// FlatMap applies a function that returns a single of V2 to the source single of V1.
// The Single from the function is flattened (hence FlatMap).
func FlatMap[A any, B any](src Single[A], apply func(A) Single[B]) Single[B] {
	return Single[B]{
		stream.FlatMap(src.ToObservable(), func(a A) stream.Observable[B] { return apply(a).ToObservable() }),
	}
}

// RetrieveValue returns the result value emitted by the Single
func RetrieveValue[T any](ctx context.Context, src Single[T]) (T, error) {
	return stream.First(ctx, src.ToObservable())
}

// MapDerefPtr takes a single of a pointer and maps it to a single of a de-referenced value of that pointer
func MapDerefPtr[T any](src Single[*T]) Single[T] {
	return fromObservable(streamx.MapDerefPtr(src.ToObservable()))
}

func fromObservable[T any](src stream.Observable[T]) Single[T] {
	return Single[T]{src: stream.Take(1, src)}
}
