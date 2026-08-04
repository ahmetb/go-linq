package linq

import (
	"reflect"
	"testing"
)

func TestGroupBy(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	// Groups are yielded in order of first appearance of their key:
	// key 1 is seen first (element 1), then key 0 (element 2).
	want := []Group[int, int]{
		{Key: 1, Group: []int{1, 3, 5, 7, 9}},
		{Key: 0, Group: []int{2, 4, 6, 8}},
	}

	q := FromSlice(input).GroupBy(
		func(i int) int { return i % 2 },
		func(i int) int { return i },
	)

	if got := toSlice(q); !reflect.DeepEqual(got, want) {
		t.Errorf("FromSlice(%v).GroupBy()=%v expected %v", input, got, want)
	}
}

func TestGroupBy_TypeChanging(t *testing.T) {
	input := []string{"apple", "avocado", "banana"}

	q := FromSlice(input).GroupBy(
		func(s string) byte { return s[0] },
		func(s string) int { return len(s) },
	)

	want := []Group[byte, int]{
		{Key: 'a', Group: []int{5, 7}},
		{Key: 'b', Group: []int{6}},
	}

	if got := toSlice(q); !reflect.DeepEqual(got, want) {
		t.Errorf("GroupBy()=%v expected %v", got, want)
	}
}

func TestGroupBy_Abort(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	q := FromSlice(input).GroupBy(
		func(i int) int { return i % 2 },
		func(i int) int { return i },
	)

	runDryIteration(q)
}
