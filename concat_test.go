package linq

import "testing"

func TestAppend(t *testing.T) {
	input := []int{1, 2, 3, 4}
	want := []int{1, 2, 3, 4, 5}

	if q := FromSlice(input).Append(5); !testQueryIteration(q, want) {
		t.Errorf("FromSlice(%v).Append()=%v expected %v", input, toSlice(q), want)
	}
}

func TestConcat(t *testing.T) {
	input1 := []int{1, 2, 3}
	input2 := []int{4, 5}
	want := []int{1, 2, 3, 4, 5}

	if q := FromSlice(input1).Concat(FromSlice(input2)); !testQueryIteration(q, want) {
		t.Errorf("FromSlice(%v).Concat(%v)=%v expected %v", input1, input2, toSlice(q), want)
	}
}

func TestPrepend(t *testing.T) {
	input := []int{1, 2, 3, 4}
	want := []int{0, 1, 2, 3, 4}

	if q := FromSlice(input).Prepend(0); !testQueryIteration(q, want) {
		t.Errorf("FromSlice(%v).Prepend()=%v expected %v", input, toSlice(q), want)
	}
}
