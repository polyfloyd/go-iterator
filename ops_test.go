package iterator

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

func TestMap(t *testing.T) {
	iter := FromSlice([]int{1, 2, 3})
	iter = Map(iter, func(i int) int { return i * 2 })
	result := ToSlice(iter)
	if !reflect.DeepEqual(result, []int{2, 4, 6}) {
		t.Fatalf("Unexpected: %v", result)
	}

	countIter := Map(FromSlice([]int{1, 2, 3}), func(i int) int { return i * 2 })
	testCounterImplementation(t, countIter, 3)
}

func TestFilterMap(t *testing.T) {
	iter := FromSlice([]int{1, 2, 3, 4})
	iter = FilterMap(iter, func(i int) (int, bool) {
		j := i * 2
		return j, j < 5
	})
	result := ToSlice(iter)
	if !reflect.DeepEqual(result, []int{2, 4}) {
		t.Fatalf("Unexpected: %v", result)
	}
}

func TestFlatten(t *testing.T) {
	iter0 := FromSlice([][]int{{0, 1, 2}, {10, 11, 12}})
	iter1 := Flatten(iter0, FromSlice[int])
	result := ToSlice(iter1)
	if !reflect.DeepEqual(result, []int{0, 1, 2, 10, 11, 12}) {
		t.Fatalf("Unexpected: %v", result)
	}

	countIter0 := Flatten(FromSlice([][]int{{0, 1, 2}, {100}, {10, 11}}), FromSlice[int])
	testCounterImplementation(t, countIter0, 6)

	countIter1 := Flatten(FromSlice([][]int{{0, 1, 2}, {100}, {10, 11}}), FromSlice[int])
	countIter1.Next() // Test whether the partially consumed iterator is included.
	testCounterImplementation(t, countIter1, 5)
}

func TestFilter(t *testing.T) {
	iter := FromSlice([]int{1, 2, 3, 4, 5, 6})
	iter = Filter(iter, func(i int) bool { return i%2 == 0 })
	result := ToSlice(iter)
	if !reflect.DeepEqual(result, []int{2, 4, 6}) {
		t.Fatalf("Unexpected: %v", result)
	}
}

func TestTake(t *testing.T) {
	testCounterImplementation(t, Take(Repeat[int](1337), 10), 10)
	testCounterImplementation(t, Take(Range[int](0, 20, 1), 10), 10)
}

// TestReduce is covered by other the other tests of the functions that use it.

func TestCount(t *testing.T) {
	iter := FromSlice([]int{0, 0, 0})
	iter = Go(context.Background(), iter) // Ensure Counter is not implemented.
	result := Count(iter)
	if result != 3 {
		t.Fatalf("Unexpected: %v", result)
	}
}

func testCounterImplementation[T any](t *testing.T, iter Iterator[T], expectedCount int) {
	t.Run(fmt.Sprintf("count %d", expectedCount), func(t *testing.T) {
		counter, ok := iter.(Counter[T])
		if !ok {
			t.Fatalf("%T does not implement Counter", iter)
		}
		if count := Count(iter); count != expectedCount {
			t.Fatalf("Unexpected count: %v", count)
		}
		if _, ok := counter.Next(); ok {
			t.Fatalf("Unexpected item after calling Count")
		}
		if count := counter.Count(); count != 0 {
			t.Fatalf("Unexpected second count: %v", count)
		}
	})
}

func TestSum(t *testing.T) {
	iter := FromSlice([]int{1, 2, 3})
	result := Sum(iter)
	if result != 6 {
		t.Fatalf("Unexpected: %v", result)
	}
}

func TestMin(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		iter := Empty[int]()
		val, ok := Min(iter)
		if ok {
			t.Fatalf("Unexpected: %v", val)
		}
	})
	t.Run("numbers", func(t *testing.T) {
		iter := FromSlice([]int{4, 5, 1, 2, 3})
		val, ok := Min(iter)
		if !ok || val != 1 {
			t.Fatalf("Unexpected: %v, %v", val, ok)
		}
	})
}

func TestMax(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		iter := Empty[int]()
		val, ok := Max(iter)
		if ok {
			t.Fatalf("Unexpected: %v", val)
		}
	})
	t.Run("numbers", func(t *testing.T) {
		iter := FromSlice([]int{4, 5, 1, 2, 3})
		val, ok := Max(iter)
		if !ok || val != 5 {
			t.Fatalf("Unexpected: %v, %v", val, ok)
		}
	})
}

func TestJoin(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		iter := Empty[string]()
		val := Join(iter, ", ")
		if val != "" {
			t.Fatalf("Unexpected: %v", val)
		}
	})
	t.Run("strings", func(t *testing.T) {
		iter := FromSlice([]string{"foo", "bar", "baz", "qux"})
		val := Join(iter, ", ")
		if val != "foo, bar, baz, qux" {
			t.Fatalf("Unexpected: %v", val)
		}
	})
}
