package iterator

import (
	"constraints"
)

// Map returns a new iterator which applies a function to all items from the input iterator which
// are subsequently returned.
//
// The mapping function should not mutate the state outside its scope.
func Map[T any, O any](from Iterator[T], mapFunc func(T) O) Iterator[O] {
	return &mapIterator[T, O]{from: from, mapFunc: mapFunc}
}

type mapIterator[T any, O any] struct {
	from Iterator[T]
	mapFunc func(T) O
}

func (iter *mapIterator[T, O]) Next() (O, bool) {
	item, ok := iter.from.Next()
	if !ok {
		var zero O
		return zero, false
	}
	mapped := iter.mapFunc(item)
	return mapped, true
}

// Filter returns a new iterator that returns only the items that pass the test of the specified
// filter function.
//
// The filter function should not mutate the state outside its scope.
func Filter[T any](from Iterator[T], filterFunc func(T) bool) Iterator[T] {
	return &filterIterator[T]{from: from, filterFunc: filterFunc}
}

type filterIterator[T any] struct {
	from Iterator[T]
	filterFunc func(T) bool
}

func (iter *filterIterator[T]) Next() (T, bool) {
	for item, ok := iter.from.Next(); ok; item, ok = iter.from.Next() {
		if iter.filterFunc(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

func Reduce[T any, O any](from Iterator[T], reduceFunc func(O, T) O, initial O) O {
	accum := initial
	for item, ok := from.Next(); ok; item, ok = from.Next() {
		accum = reduceFunc(accum, item)
	}
	return accum
}

// Sum adds all the items from the iterator.
func Sum[T Number](from Iterator[T]) T {
	var zero T
	return Reduce(from, func(accum T, item T) T {
		return accum + item
	}, zero)
}

func Min[T constraints.Ordered](from Iterator[T]) (T, bool) {
	init, ok := from.Next()
	if !ok {
		var zero T
		return zero, false
	}
	min := Reduce(from, func(accum T, item T) T {
		if item < accum {
			return item
		}
		return accum
	}, init)
	return min, true
}

func Max[T constraints.Ordered](from Iterator[T]) (T, bool) {
	init, ok := from.Next()
	if !ok {
		var zero T
		return zero, false
	}
	max := Reduce(from, func(accum T, item T) T {
		if item > accum {
			return item
		}
		return accum
	}, init)
	return max, true
}

// Join concatenates the strings from an iterator into a single string, with the items separated by
// the specified separator string.
func Join[T ~string](from Iterator[T], sep string) string {
	following := false
	accum := ""
	for item, ok := from.Next(); ok; item, ok = from.Next() {
		if following {
			accum += sep
		}
		following = true
		accum += string(item)
	}
	return accum
}
