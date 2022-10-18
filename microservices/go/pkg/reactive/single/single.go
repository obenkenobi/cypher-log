package single

import (
	"context"
	"github.com/akrennmair/slice"
	"github.com/barweiss/go-tuple"
	"github.com/joamaki/goreactive/stream"
	"sync"
)

// Single is a listener for a single value. It by default listens for a value
// lazily or if specified, asynchronously. Async singles will cache the stored
// result.
type Single[T any] struct {
	src stream.Observable[T]
}

// ToObservable turns the Single into an observable. (Warning: the provided
// single can no longer reliably be used to read values as they may be used up by
// the observable instead.)
func (s Single[T]) ToObservable() stream.Observable[T] { return s.src }

// ScheduleEagerAsyncCached takes a single and returns a new Single that is scheduled
// to be evaluated eagerly and asynchronously up until point of execution of the
// returning single. The emitted value is cached if observed twice.
func (s Single[T]) ScheduleEagerAsyncCached(ctx context.Context) Single[T] {
	ch, errCh := ToChannels(ctx, s)
	return FromChannelsCached(ch, errCh)
}

// ScheduleLazyAndCache takes a single and returns a new single that is lazy
// evaluated. The single's emitted value is cached if observed twice.
func (s Single[T]) ScheduleLazyAndCache(ctx context.Context) Single[T] {
	return FromSupplierCached(func() (T, error) {
		return RetrieveValue(ctx, s)
	})
}

// Thread safe Observable that reads from a channel, used for ensuring when a
// single is made from a channel, it will always emit the same value even when
// observed multiple times. This is done by caching the channel result in a
// thread safe manner.
type singleChanReadObservable[T any] struct {
	_channelRead bool
	_ch          <-chan T
	_chErr       <-chan error
	_value       T
	_error       error
	_valueRWLock sync.RWMutex
}

func (a *singleChanReadObservable[T]) Observe(ctx context.Context, next func(T) error) error {
	if ctx.Err() != nil {
		// Context already cancelled, stop before emitting items.
		return ctx.Err()
	}
	shouldAttemptWriteValue := func() bool {
		a._valueRWLock.RLock()
		defer a._valueRWLock.RUnlock()
		return !a._channelRead
	}()

	if shouldAttemptWriteValue {
		func() {
			a._valueRWLock.Lock()
			defer a._valueRWLock.Unlock()
			if !a._channelRead {
				a._channelRead = true
				if a._chErr != nil {
					if err := <-a._chErr; err != nil {
						a._error = err
					}
				}
				if a._error == nil {
					a._value = <-a._ch
				}
			}
		}()
	}
	value, err := func() (T, error) {
		a._valueRWLock.RLock()
		defer a._valueRWLock.RUnlock()
		return a._value, a._error
	}()
	if err != nil {
		return err
	}
	return next(value)
}

type singleSupplierReadObservable[T any] struct {
	_supplierRan bool
	_supplier    func() (T, error)
	_value       T
	_err         error
	_valueRWLock sync.RWMutex
}

func (a *singleSupplierReadObservable[T]) Observe(ctx context.Context, next func(T) error) error {
	if ctx.Err() != nil {
		// Context already cancelled, stop before emitting items.
		return ctx.Err()
	}
	shouldAttemptWriteValue := func() bool {
		a._valueRWLock.RLock()
		defer a._valueRWLock.RUnlock()
		return !a._supplierRan
	}()

	if shouldAttemptWriteValue {
		func() {
			a._valueRWLock.Lock()
			defer a._valueRWLock.Unlock()
			if !a._supplierRan {
				val, err := a._supplier()
				if err != nil {
					a._err = err
				} else {
					a._value = val
				}
				a._supplierRan = true
			}
		}()
	}
	value, err := func() (T, error) {
		a._valueRWLock.RLock()
		defer a._valueRWLock.RUnlock()
		return a._value, a._err
	}()
	if err != nil {
		return err
	}
	return next(value)
}

func FromObservableAsList[T any](obs stream.Observable[T]) Single[[]T] {
	listObs := stream.Reduce(obs, []T{}, func(list []T, v T) []T { return append(list, v) })
	return fromSingleObservable(listObs)
}

// FromChannelCached Creates a single that listens to a single value from a channel.
// Recommended for channels that only emmit one value. Once the channel is read,
// the value is cached so the single can be observed more than one time with
// little overhead. If a channel emits multiple values, it is recommended you use
// observables instead.
func FromChannelCached[T any](ch <-chan T) Single[T] { return FromChannelsCached(ch, nil) }

// FromChannelsCached Creates a single that listens to a single value from a channel
// and checks for errors in the error channel. Once the channel is read, the
// value is cached so the single can be observed more than one time with little
// overhead. Recommended for channels that only emmit one value. If a channel
// emits multiple values, it is recommended you use observables instead.
func FromChannelsCached[T any](ch <-chan T, chErr <-chan error) Single[T] {
	return Single[T]{src: &singleChanReadObservable[T]{_ch: ch, _chErr: chErr}}
}

// Just creates a single from a single item
func Just[T any](val T) Single[T] { return fromSingleObservable(stream.Just(val)) }

// Error creates a single that fails immediately with the given error
func Error[T any](err error) Single[T] { return fromSingleObservable(stream.Error[T](err)) }

// FromSupplierCached creates a single out of a supplier function that returns a
// value or an error to emit. If the supplier is successful, the result value is
// emitted. Otherwise, if an error is returned, an error is emitted. When this
// single is observed, the emitted value from the supplier is cached if observed
// again.
func FromSupplierCached[T any](supplier func() (T, error)) Single[T] {
	var src stream.Observable[T] = &singleSupplierReadObservable[T]{_supplier: supplier}
	return fromSingleObservable(src)
}

// Map applies a function onto a Single
func Map[A any, B any](src Single[A], apply func(A) B) Single[B] {
	return fromSingleObservable[B](stream.Map(src.ToObservable(), func(a A) B { return apply(a) }))
}

// MapWithError applies a function onto a Single where if an error is returned, the Single fails
func MapWithError[A any, B any](src Single[A], apply func(A) (B, error)) Single[B] {
	return fromSingleObservable[B](
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
	return FlatMap(src1, func(v1 V1) Single[tuple.T2[V1, V2]] {
		return Map(src2, func(v2 V2) tuple.T2[V1, V2] {
			return tuple.New2(v1, v2)
		})

	})
}

// Zip3 Takes 3 Singles and returns a Single that emits a tuple of each of the
// singles in the order they are supplied
func Zip3[V1 any, V2 any, V3 any](
	src1 Single[V1],
	src2 Single[V2],
	src3 Single[V3],
) Single[tuple.T3[V1, V2, V3]] {
	return FlatMap(
		Zip2(src1, src2),
		func(t tuple.T2[V1, V2]) Single[tuple.T3[V1, V2, V3]] {
			return Map(src3, func(v3 V3) tuple.T3[V1, V2, V3] {
				return tuple.New3(t.V1, t.V2, v3)
			})
		},
	)
}

// Zip4 Takes 4 Singles and returns a Single that emits a tuple of each of the
// singles in the order they are supplied.
func Zip4[V1 any, V2 any, V3 any, V4 any](
	src1 Single[V1],
	src2 Single[V2],
	src3 Single[V3],
	src4 Single[V4],
) Single[tuple.T4[V1, V2, V3, V4]] {
	return FlatMap(
		Zip2(src1, src2),
		func(t1 tuple.T2[V1, V2]) Single[tuple.T4[V1, V2, V3, V4]] {
			return Map(Zip2(src3, src4), func(t2 tuple.T2[V3, V4]) tuple.T4[V1, V2, V3, V4] {
				return tuple.New4(t1.V1, t1.V2, t2.V1, t2.V2)
			})
		},
	)
}

// Zip5 Takes 5 Singles and returns a Single that emits a tuple of each of the
// singles in the order they are supplied.
func Zip5[V1 any, V2 any, V3 any, V4 any, V5 any](
	src1 Single[V1],
	src2 Single[V2],
	src3 Single[V3],
	src4 Single[V4],
	src5 Single[V5],
) Single[tuple.T5[V1, V2, V3, V4, V5]] {
	return FlatMap(
		Zip3(src1, src2, src3),
		func(t1 tuple.T3[V1, V2, V3]) Single[tuple.T5[V1, V2, V3, V4, V5]] {
			return Map(Zip2(src4, src5), func(t2 tuple.T2[V4, V5]) tuple.T5[V1, V2, V3, V4, V5] {
				return tuple.New5(t1.V1, t1.V2, t1.V3, t2.V1, t2.V2)
			})
		},
	)
}

// Zip6 Takes 6 Singles and returns a Single that emits a tuple of each of the
// singles in the order they are supplied.
func Zip6[V1 any, V2 any, V3 any, V4 any, V5 any, V6 any](
	src1 Single[V1],
	src2 Single[V2],
	src3 Single[V3],
	src4 Single[V4],
	src5 Single[V5],
	src6 Single[V6],
) Single[tuple.T6[V1, V2, V3, V4, V5, V6]] {
	return FlatMap(
		Zip3(src1, src2, src3),
		func(t1 tuple.T3[V1, V2, V3]) Single[tuple.T6[V1, V2, V3, V4, V5, V6]] {
			return Map(Zip3(src4, src5, src6), func(t2 tuple.T3[V4, V5, V6]) tuple.T6[V1, V2, V3, V4, V5, V6] {
				return tuple.New6(t1.V1, t1.V2, t1.V3, t2.V1, t2.V2, t2.V3)
			})
		},
	)
}

func MergeSingles[T any](singles []Single[T]) stream.Observable[T] {
	observables := slice.Map(singles, func(s Single[T]) stream.Observable[T] {
		return s.ToObservable()
	})
	return stream.Merge(observables...)
}

// FlatMap applies a function that returns a single of V2 to the source single of V1.
// The Single from the function is flattened (hence FlatMap).
func FlatMap[A any, B any](src Single[A], apply func(A) Single[B]) Single[B] {
	return Single[B]{
		stream.FlatMap(src.ToObservable(), func(a A) stream.Observable[B] { return apply(a).ToObservable() }),
	}
}

// RetrieveValue returns the value emitted by the Single
func RetrieveValue[T any](ctx context.Context, src Single[T]) (T, error) {
	return stream.First(ctx, src.ToObservable())
}

// ToChannels returns channels that can emit a value created by a single. It
// ensures asynchronous execution when the value is evaluated.
func ToChannels[T any](ctx context.Context, src Single[T]) (<-chan T, <-chan error) {
	return stream.ToChannels(ctx, src.ToObservable())
}

func fromSingleObservable[T any](src stream.Observable[T]) Single[T] {
	return Single[T]{src: src}
}
