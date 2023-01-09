package concurrent

import (
	"github.com/barweiss/go-tuple"
	"sync"
)

type Future[T any] interface {
	Await() (T, error)
}

type FutureImpl[T any] struct {
	_channelRead bool
	_ch          <-chan tuple.T2[T, error]
	_value       T
	_error       error
	_valueRWLock sync.RWMutex
}

func (a FutureImpl[T]) Await() (T, error) {
	shouldReadChannel := func() bool {
		a._valueRWLock.RLock()
		defer a._valueRWLock.RUnlock()
		return !a._channelRead
	}()

	if shouldReadChannel {
		func() {
			a._valueRWLock.Lock()
			defer a._valueRWLock.Unlock()
			if !a._channelRead {
				a._channelRead = true
				t := <-a._ch
				a._value = t.V1
				a._error = t.V2
			}
		}()
	}
	return func() (T, error) {
		a._valueRWLock.RLock()
		defer a._valueRWLock.RUnlock()
		return a._value, a._error
	}()
}

func Async[T any](supplier func() (T, error)) *FutureImpl[T] {
	ch := make(chan tuple.T2[T, error])
	go func() {
		defer close(ch)
		v, err := supplier()
		ch <- tuple.New2(v, err)
	}()
	return &FutureImpl[T]{_channelRead: false, _error: nil, _ch: ch}
}
