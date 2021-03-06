package iterator

import (
	"golang.org/x/exp/constraints"
)

// Map returns a new iterator which applies a function to all items from the input iterator which
// are subsequently returned.
//
// The mapping function should not mutate the state outside its scope.
func Map[T any, O any](from Iterator[T], mapFunc func(T) O) Iterator[O] {
	return &mapIterator[T, O]{from: from, mapFunc: mapFunc}
}

type mapIterator[T any, O any] struct {
	from    Iterator[T]
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

func (iter *mapIterator[T, O]) Count() int {
	return Count(iter.from)
}

// FilterMap applies a function to all items from the specified iterator as Map does, but culls the
// results which are accompanied by false.
func FilterMap[T any, O any](from Iterator[T], mapFunc func(T) (O, bool)) Iterator[O] {
	return &filterMapIterator[T, O]{from: from, mapFunc: mapFunc}
}

type filterMapIterator[T any, O any] struct {
	from    Iterator[T]
	mapFunc func(T) (O, bool)
}

func (iter *filterMapIterator[T, O]) Next() (O, bool) {
	for item, ok := iter.from.Next(); ok; item, ok = iter.from.Next() {
		mapped, ok := iter.mapFunc(item)
		if ok {
			return mapped, true
		}
	}
	var zero O
	return zero, false
}

// Flatten applies a function to all items of the specified iterator, returning an iterator for each
// item. The resulting iterators are then concatenated into a single iterator.
func Flatten[T any](from Iterator[Iterator[T]]) Iterator[T] {
	return &flattenIterator[T]{from: from}
}

type flattenIterator[T any] struct {
	from Iterator[Iterator[T]]
	head Iterator[T]
}

func (iter *flattenIterator[T]) Next() (T, bool) {
	for {
		if iter.head == nil {
			item, ok := iter.from.Next()
			if !ok {
				var zero T
				return zero, false
			}
			iter.head = item
		}
		item, ok := iter.head.Next()
		if ok {
			return item, true
		}
		iter.head = nil
	}
}

func (iter *flattenIterator[T]) Count() int {
	fromCounts := Map(iter.from, func(item Iterator[T]) int {
		return Count(item)
	})
	count := Sum(fromCounts)
	if iter.head != nil {
		count += Count(iter.head)
	}
	return count
}

// Filter returns a new iterator that returns only the items that pass the test of the specified
// filter function.
//
// The filter function should not mutate the state outside its scope.
func Filter[T any](from Iterator[T], filterFunc func(T) bool) Iterator[T] {
	return &filterIterator[T]{from: from, filterFunc: filterFunc}
}

type filterIterator[T any] struct {
	from       Iterator[T]
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

// Take limits the number of items returned by an iterator to the specified count.
func Take[T any](from Iterator[T], num int) Iterator[T] {
	return &takeIterator[T]{from: from, num: num}
}

type takeIterator[T any] struct {
	from Iterator[T]
	num  int
}

func (iter *takeIterator[T]) Next() (T, bool) {
	if iter.num <= 0 {
		var zero T
		return zero, false
	}
	item, ok := iter.from.Next()
	if ok {
		iter.num--
	}
	return item, ok
}

func (iter *takeIterator[T]) Count() int {
	count := 0
	if counter, ok := iter.from.(Counter[T]); ok {
		count = counter.Count()
	} else {
		// If `Count(from)` is larger than num we are potentially doing more work than needed. So we
		// have our own count loop with a limit on the loop condition that exits early.
		for _, ok := iter.from.Next(); ok && count < iter.num; _, ok = iter.from.Next() {
			count++
		}
	}

	if iter.num < count {
		count = iter.num
	}
	iter.num = 0
	return count
}

func Reduce[T any, O any](from Iterator[T], reduceFunc func(O, T) O, initial O) O {
	accum := initial
	for item, ok := from.Next(); ok; item, ok = from.Next() {
		accum = reduceFunc(accum, item)
	}
	return accum
}

// Counter can optionally be implemented by iterators to provide a specialized implementation of
// Count. Implementations must ensure that after Count was called, Next will return no more items.
type Counter[T any] interface {
	Iterator[T]
	Count() int
}

// Count consumes the entire iterator and returns the number of remaining elements that were
// returned.
func Count[T any](from Iterator[T]) int {
	if counter, ok := from.(Counter[T]); ok {
		return counter.Count()
	}

	count := 0
	for _, ok := from.Next(); ok; _, ok = from.Next() {
		count++
	}
	return count
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
