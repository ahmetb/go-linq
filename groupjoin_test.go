package linq

import "testing"

func TestGroupJoin(t *testing.T) {
	outer := []int{0, 1, 2}
	inner := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	want := []KeyValue[int, int]{
		{0, 4},
		{1, 5},
		{2, 0},
	}

	q := FromSlice(outer).GroupJoin(
		FromSlice(inner),
		func(i int) int { return i },
		func(i int) int { return i % 2 },
		func(outer int, inners []int) KeyValue[int, int] {
			return KeyValue[int, int]{outer, len(inners)}
		})

	if !testQueryIteration(q, want) {
		t.Errorf("FromSlice().GroupJoin()=%v expected %v", toSlice(q), want)
	}
}
