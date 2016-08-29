package linq

import "testing"

func TestAppend(t *testing.T) {
	input := []int{1, 2, 3, 4}
	want := []interface{}{1, 2, 3, 4, 5}

	if q := From(input).Append(5); !validateQuery(q, want) {
		t.Errorf("From(%v).Append()=%v expected %v", input, toSlice(q), want)
	}
}

func TestConcat(t *testing.T) {
	input1 := []int{1, 2, 3}
	input2 := []int{4, 5}
	want := []interface{}{1, 2, 3, 4, 5}

	if q := From(input1).Concat(From(input2)); !validateQuery(q, want) {
		t.Errorf("From(%v).Concat(%v)=%v expected %v", input1, input2, toSlice(q), want)
	}
}

func TestPrepend(t *testing.T) {
	input := []int{1, 2, 3, 4}
	want := []interface{}{0, 1, 2, 3, 4}

	if q := From(input).Prepend(0); !validateQuery(q, want) {
		t.Errorf("From(%v).Prepend()=%v expected %v", input, toSlice(q), want)
	}
}
