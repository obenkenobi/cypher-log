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

// ScheduleEagerAsync takes a single and returns a new Single that is scheduled
// to be evaluated eagerly and asynchronously up until point of execution of the
// returning single.
func (s Single[T]) ScheduleEagerAsync(ctx context.Context) Single[T] {
	ch, errCh := ToChannels(ctx, s)
	return FromChannels(ch, errCh)
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
	_valueRWLock sync.RWMutex
}

func (a *singleChanReadObservable[T]) Observe(ctx context.Context, next func(T) error) error {
	if ctx.Err() != nil {
		// Context already cancelled, stop before emitting items.
		return ctx.Err()
	}
	a._valueRWLock.RLock()
	shouldAttemptWriteValue := !a._channelRead
	a._valueRWLock.RUnlock()

	if shouldAttemptWriteValue {
		a._valueRWLock.Lock()
		if !a._channelRead {
			a._channelRead = true
			if a._chErr != nil {
				if err := <-a._chErr; err != nil {
					return err
				}
			}
			a._value = <-a._ch
		}
		a._valueRWLock.Unlock()
	}
	a._valueRWLock.RLock()
	value := a._value
	a._valueRWLock.RUnlock()
	return next(value)
}

func FromObservableAsList[T any](obs stream.Observable[T]) Single[[]T] {
	listObs := stream.Reduce(obs, []T{}, func(res []T, v T) []T { return append(res, v) })
	return fromSingleObservable(listObs)
}

// FromChannel Creates a single that listens to a single value from a channel.
// Recommended for channels that only emmit one value. If a channel emits
// multiple values, it is recommended you use observables instead.
func FromChannel[T any](ch <-chan T) Single[T] { return FromChannels(ch, nil) }

// FromChannels Creates a single that listens to a single value from a channel
// and checks for errors in the error channel. Recommended for channels that only
// emmit one value. If a channel emits multiple values, it is recommended you use
// observables instead.
func FromChannels[T any](ch <-chan T, chErr <-chan error) Single[T] {
	return Single[T]{src: &singleChanReadObservable[T]{_ch: ch, _chErr: chErr}}
}

// Just creates a single from a single item
func Just[T any](val T) Single[T] { return fromSingleObservable(stream.Just(val)) }

// Error creates a single that fails immediately with the given error
func Error[T any](err error) Single[T] { return fromSingleObservable(stream.Error[T](err)) }

// FromSupplier
//creates a single out of a supplier function that returns a value or an error
//to emit. If the supplier is successful, the result value is emitted.
//Otherwise, if an error is returned, an error id emitted.
func FromSupplier[T any](supplier func() (T, error)) Single[T] {
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
