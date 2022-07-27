package reactorextensions

import (
	"context"
	"github.com/joamaki/goreactive/stream"
)

func ObserveProducer[V any](producer func() (V, error)) stream.Observable[V] {
	return stream.FuncObservable[V](func(ctx context.Context, next func(V) error) error {
		if ctx.Err() != nil {
			// Context already cancelled, stop before emitting items.
			return ctx.Err()
		}
		val, err := producer()
		if err != nil {
			return err
		}
		return next(val)
	})
}
