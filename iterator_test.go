package iterator

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"golang.org/x/exp/constraints"
)

func TestEmpty(t *testing.T) {
	t.Run("next", func(t *testing.T) {
		iter := Empty[int]()
		val, ok := iter.Next()
		if ok {
			t.Fatalf("Unexpected: %v", val)
		}
	})

	testCounterImplementation(t, Empty[int](), 0)
}

func TestOnce(t *testing.T) {
	t.Run("next", func(t *testing.T) {
		iter := Once[int](1337)
		result := ToSlice(iter)
		if !reflect.DeepEqual(result, []int{1337}) {
			t.Fatalf("Unexpected: %v", result)
		}
	})

	testCounterImplementation(t, Once[int](1337), 1)
}

func TestRepeat(t *testing.T) {
	iter := Repeat[int](1337)
	iter = Take(iter, 4)
	result := ToSlice(iter)
	if !reflect.DeepEqual(result, []int{1337, 1337, 1337, 1337}) {
		t.Fatalf("Unexpected: %v", result)
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

	testCounterImplementation(t, Range[int](0, 10, 1), 10)
	testCounterImplementation(t, Range[int](0, 10, 2), 5)
	testCounterImplementation(t, Range[int](0, 10, 5), 2)
	testCounterImplementation(t, Range[int](0, 0, 1), 0)
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

	testCounterImplementation(t, FromSlice[int]([]int{}), 0)
	testCounterImplementation(t, FromSlice[int]([]int{1, 2, 3, 4}), 4)
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

func TestFromMap(t *testing.T) {
	m := map[string]int{
		"x": 1,
		"y": 2,
		"z": 3,
	}
	iter := FromMap(m)
	result := ToSlice(iter)
	sort.Sort(byMapEntryKey[string, int](result)) // Map iteration is non-deterministic.
	expect := []MapEntry[string, int]{
		{Key: "x", Val: 1},
		{Key: "y", Val: 2},
		{Key: "z", Val: 3},
	}
	if !reflect.DeepEqual(result, expect) {
		t.Fatalf("Unexpected: %v", result)
	}
}

func TestToMap(t *testing.T) {
	iter := FromSlice([]MapEntry[string, int]{
		{Key: "x", Val: 1},
		{Key: "y", Val: 2},
		{Key: "z", Val: 3},
	})
	result := ToMap(iter)
	expect := map[string]int{
		"x": 1,
		"y": 2,
		"z": 3,
	}
	if !reflect.DeepEqual(result, expect) {
		t.Fatalf("Unexpected: %v", result)
	}
}

type byMapEntryKey[K constraints.Ordered, V any] []MapEntry[K, V]

func (s byMapEntryKey[K, V]) Len() int           { return len(s) }
func (s byMapEntryKey[K, V]) Less(i, j int) bool { return s[i].Key < s[j].Key }
func (s byMapEntryKey[K, V]) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
