package linq

import "testing"

func TestIntersect(t *testing.T) {
	input1 := []int{1, 2, 3}
	input2 := []int{1, 4, 7, 9, 12, 3}
	want := []int{1, 3}

	if q := FromSlice(input1).Intersect(FromSlice(input2)); !testQueryIteration(q, want) {
		t.Errorf("FromSlice(%v).Intersect(%v)=%v expected %v", input1, input2, toSlice(q), want)
	}
}

func TestIntersectBy(t *testing.T) {
	input1 := []int{5, 7, 8}
	input2 := []int{1, 4, 7, 9, 12, 3}
	want := []int{5, 8}

	if q := FromSlice(input1).IntersectBy(FromSlice(input2), func(i int) int {
		return i % 2
	}); !testQueryIteration(q, want) {
		t.Errorf("FromSlice(%v).IntersectBy(%v)=%v expected %v", input1, input2, toSlice(q), want)
	}
}
