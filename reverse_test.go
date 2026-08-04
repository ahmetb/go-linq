package linq

import "testing"

func TestReverse(t *testing.T) {
	input := []int{1, 2, 3}
	want := []int{3, 2, 1}

	if q := FromSlice(input).Reverse(); !testQueryIteration(q, want) {
		t.Errorf("FromSlice(%v).Reverse()=%v expected %v", input, toSlice(q), want)
	}
}
