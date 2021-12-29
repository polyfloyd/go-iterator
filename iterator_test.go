package iterator

import (
	"context"
	"reflect"
	"testing"
)

func TestEmpty(t *testing.T) {
	iter := Empty[int]()
	val, ok := iter.Next()
	if ok {
		t.Fatalf("Unexpected: %v", val)
	}
}

func TestRange(t *testing.T) {
	t.Run("count to 5", func(t *testing.T) {
		iter := Range[int](0, 5, 1)
		result := ToSlice(iter)
		if !reflect.DeepEqual(result, []int{0, 1, 2, 3, 4}) {
			t.Fatalf("Unexpected: %v", result)
		}
	})
	t.Run("empty", func(t *testing.T) {
		iter := Range[int](0, 0, 1)
		result := ToSlice(iter)
		if !reflect.DeepEqual(result, []int{}) {
			t.Fatalf("Unexpected: %v", result)
		}
	})
	t.Run("panic on end before start", func(t *testing.T) {
		var err interface{}
		func() {
			defer func() { err = recover() }()
			Range[int](4, 0, 1)
		}()
		if err == nil {
			t.Fatalf("Expected panic")
		}
	})
	t.Run("panic on zero step", func(t *testing.T) {
		var err interface{}
		func() {
			defer func() { err = recover() }()
			Range[int](0, 4, 0)
		}()
		if err == nil {
			t.Fatalf("Expected panic")
		}
	})
}

func TestFromSlice(t *testing.T) {
	t.Run("items", func(t *testing.T) {
		iter := FromSlice[int]([]int{1, 2, 3, 4})
		result := ToSlice[int](iter)
		if !reflect.DeepEqual(result, []int{1, 2, 3, 4}) {
			t.Fatalf("Unexpected: %v", result)
		}
	})
	t.Run("empty", func(t *testing.T) {
		iter := FromSlice[int]([]int{})
		result := ToSlice[int](iter)
		if !reflect.DeepEqual(result, []int{}) {
			t.Fatalf("Unexpected: %v", result)
		}
	})
}

// ToSlice is already quite well covered because it is used in other tests.

func TestToChannel(t *testing.T) {
	t.Run("cancel unconsumed", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		iter := FromSlice([]int{1, 2, 3, 4})
		iter = FromChannel(ToChannel(ctx, iter, 0)) // Unbuffered, to avoid a race condition
		iter.Next()
		cancel()
		val, ok := iter.Next()
		if ok {
			t.Fatalf("Unexpected: %v", val)
		}
	})
}

func TestGo(t *testing.T) {
	t.Run("items", func(t *testing.T) {
		iter := FromSlice([]int{1, 2, 3, 4})
		iter = Go(context.Background(), iter)
		result := ToSlice[int](iter)
		if !reflect.DeepEqual(result, []int{1, 2, 3, 4}) {
			t.Fatalf("Unexpected: %v", result)
		}
	})
}
