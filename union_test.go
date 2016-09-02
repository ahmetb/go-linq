package linq

import "testing"

func TestUnion(t *testing.T) {
	input1 := []int{1, 2, 3}
	input2 := []int{2, 4, 5, 1}
	want := []interface{}{1, 2, 3, 4, 5}

	if q := From(input1).Union(From(input2)); !validateQuery(q, want) {
		t.Errorf("From(%v).Union(%v)=%v expected %v", input1, input2, toSlice(q), want)
	}
}
