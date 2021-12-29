package iterator

import (
	"context"
)

// An Iterator is a stream of items of some type.
//
// New instances may be created from slices, channels, integer ranges or values using the respective
// functions provided by this library.
type Iterator[T any] interface {
	// Next fetches the next item from the iterator along with an ok bit to indicate whether the
	// returned value must be considered valid.
	//
	// Whenever false is returned, it is valid to call Next again. In this case Next will keep
	// returning false indefinitely.
	Next() (T, bool)
}

// Empty returns an iterator that never returns anything.
func Empty[T any]() Iterator[T] {
	return emptyIterator[T]{}
}

type emptyIterator[T any] struct{}

func (emptyIterator[T]) Next() (T, bool) {
	var zero T
	return zero, false
}

// Range creates an iterator which returns the numeric range between start inclusive and end
// exclusive by the step size.
//
// If any of the constraints below are not met, Range will panic:
// * start <= end
// * 0 < step
func Range[T Number](start, end, step T) Iterator[T] {
	if end < start {
		panic("Range: end may not be before start")
	} else if step <= 0 {
		panic("Range: step may not be 0 or negative")
	}
	return &rangeIterator[T]{start, end, step}
}

type rangeIterator[T Number] struct {
	start, end, step T
}

func (iter *rangeIterator[T]) Next() (T, bool) {
	if iter.start >= iter.end {
		var zero T
		return zero, false
	}
	num := iter.start
	iter.start += iter.step
	return num, true
}

// FromSlice creates a new iterator which returns all items from the slice starting at index 0 until
// all items are consumed.
func FromSlice[T any](slice []T) Iterator[T] {
	return &sliceIterator[T]{slice: slice}
}

type sliceIterator[T any] struct {
	slice []T
}

func (iter *sliceIterator[T]) Next() (T, bool) {
	if len(iter.slice) == 0 {
		var zero T
		return zero, false
	}
	item := iter.slice[0]
	iter.slice = iter.slice[1:]
	return item, true
}

// ToSlice collects the items from the specified iterator into a slice.
func ToSlice[T any](from Iterator[T]) []T {
	slice := []T{}
	for item, ok := from.Next(); ok; item, ok = from.Next() {
		slice = append(slice, item)
	}
	return slice
}

func FromChannel[T any](from <-chan T) Iterator[T] {
	return &channelIterator[T]{from: from}
}

type channelIterator[T any] struct {
	from <-chan T
}

func (iter *channelIterator[T]) Next() (T, bool) {
	item, ok := <-iter.from
	return item, ok
}

// ToChannel spawns a new goroutine that pulls from the specified iterator into the returned
// channel. The channel may be buffered, which causes the preceding iterator chain to run in
// parallel to the routine that consumes from the channel.
//
// A valid context should be passed that cancels when the iterator chain goes out of scope, this
// prevents the goroutine from leaking if the channel is not fully consumed.
func ToChannel[T any](ctx context.Context, from Iterator[T], buffer int) <-chan T {
	out := make(chan T, buffer)
	go func() {
		defer close(out)
		for item, ok := from.Next(); ok; item, ok = from.Next() {
			select {
			case out <- item:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// Go is a convenience function that calls ToChannel and then FromChannel with a buffer size of 1.
//
// The effect of this is that the iterator chain preceding this call runs in parallel to
// subsequent chains.
//
// A valid context should be passed that cancels when the iterator chain goes out of scope, this
// prevents the goroutine from leaking if the iterator chain is not fully consumed.
func Go[T any](ctx context.Context, from Iterator[T]) Iterator[T] {
	return FromChannel(ToChannel(ctx, from, 1))
}
