package option

import (
	"reflect"
)

type (
	Maybe[V any] interface {
		IsPresent() bool
		IsEmpty() bool
		IfPresent(func(V))
		Filter(func(V) bool) Maybe[V]
		OrElse(other V) V
		OrElseGet(func() V) V
		Get() (V, bool)
	}
)

type some[V any] struct {
	value V
}

func (s some[V]) IsPresent() bool {
	return !s.IsEmpty()
}

func (s some[V]) IfPresent(action func(V)) {
	if s.IsPresent() {
		action(s.value)
	}
}

func (s some[V]) IsEmpty() bool {
	return isNilPtr(s.value)
}

func (s some[V]) Filter(predicate func(V) bool) Maybe[V] {
	if s.IsPresent() && predicate(s.value) {
		return s
	}
	return None[V]()
}

func (s some[V]) OrElseGet(other func() V) V {
	if s.IsEmpty() {
		return other()
	}
	return s.value
}

func (s some[V]) OrElse(other V) V {
	if s.IsEmpty() {
		return other
	}
	return s.value
}

func (s some[V]) Get() (V, bool) {
	if s.IsEmpty() {
		var s V
		return s, false
	}
	return s.value, true
}

type none[V any] struct {
}

func (n none[V]) IsPresent() bool {
	return false
}

func (n none[V]) IsEmpty() bool {
	return true
}

func (n none[V]) IfPresent(_ func(V)) {
}

func (n none[V]) Filter(_ func(V) bool) Maybe[V] {
	return n
}

func (n none[V]) OrElseGet(other func() V) V {
	return other()
}

func (n none[V]) OrElse(other V) V {
	return other
}

func (n none[V]) Get() (V, bool) {
	var val V
	return val, false
}

func Perhaps[V any](value V) Maybe[V] {
	if isNilPtr(value) {
		return None[V]()
	}
	return some[V]{value: value}
}

func None[V any]() Maybe[V] {
	return none[V]{}
}

func FlatMap[V1 any, V2 any](m Maybe[V1], mapper func(V1) Maybe[V2]) Maybe[V2] {
	if m.IsEmpty() {
		return None[V2]()
	}
	if res, ok := m.Get(); !ok || isNilPtr(res) {
		return None[V2]()
	} else {
		return mapper(res)
	}
}

func Map[V1 any, V2 any](m Maybe[V1], mapper func(V1) V2) Maybe[V2] {
	return FlatMap(m, func(v V1) Maybe[V2] {
		return Perhaps(mapper(v))
	})
}

func isNilPtr[V any](item V) bool {
	value := reflect.ValueOf(item)
	return value.IsZero() && value.Kind() == reflect.Ptr
}
