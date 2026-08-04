package linq

import "testing"

func TestExcept(t *testing.T) {
	input1 := []int{1, 2, 3, 4, 5, 1, 2, 5}
	input2 := []int{1, 2}
	want := []int{3, 4, 5, 5}

	if q := FromSlice(input1).Except(FromSlice(input2)); !testQueryIteration(q, want) {
		t.Errorf("FromSlice(%v).Except(%v)=%v expected %v", input1, input2, toSlice(q), want)
	}
}

func TestExceptBy(t *testing.T) {
	input1 := []int{1, 2, 3, 4, 5, 1, 2, 5}
	input2 := []int{1}
	want := []int{2, 4, 2}

	if q := FromSlice(input1).ExceptBy(FromSlice(input2), func(i int) int {
		return i % 2
	}); !testQueryIteration(q, want) {
		t.Errorf("FromSlice(%v).ExceptBy(%v)=%v expected %v", input1, input2, toSlice(q), want)
	}
}
