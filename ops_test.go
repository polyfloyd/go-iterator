package iterator

import (
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
	iter1 := Flatten(iter0, func(ii []int) Iterator[int] {
		return FromSlice(ii)
	})
	result := ToSlice(iter1)
	if !reflect.DeepEqual(result, []int{0, 1, 2, 10, 11, 12}) {
		t.Fatalf("Unexpected: %v", result)
	}
}

func TestFilter(t *testing.T) {
	iter := FromSlice([]int{1, 2, 3, 4, 5, 6})
	iter = Filter(iter, func(i int) bool { return i%2 == 0 })
	result := ToSlice(iter)
	if !reflect.DeepEqual(result, []int{2, 4, 6}) {
		t.Fatalf("Unexpected: %v", result)
	}
}

// TestReduce is covered by other the other tests of the functions that use it.

func TestCount(t *testing.T) {
	iter := FromSlice([]int{0, 0, 0})
	result := Count(iter)
	if result != 3 {
		t.Fatalf("Unexpected: %v", result)
	}
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
