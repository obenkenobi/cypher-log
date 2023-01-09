package concurrent

import "sync"

type Future[T any] interface {
	Await() (T, error)
}

type FutureImpl[T any] struct {
	_channelRead bool
	_ch          <-chan T
	_chErr       <-chan error
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
	return func() (T, error) {
		a._valueRWLock.RLock()
		defer a._valueRWLock.RUnlock()
		return a._value, a._error
	}()
}

func Async[T any](supplier func() (T, error)) *FutureImpl[T] {
	ch := make(chan T)
	chErr := make(chan error)
	go func() {
		defer func() {
			close(ch)
			close(chErr)
		}()
		v, err := supplier()
		ch <- v
		chErr <- err
	}()
	return &FutureImpl[T]{_channelRead: false, _error: nil}
}
