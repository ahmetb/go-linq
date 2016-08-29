package linq

import "testing"

func TestGroupJoin(t *testing.T) {
	outer := []int{0, 1, 2}
	inner := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	want := []interface{}{
		KeyValue{0, 4},
		KeyValue{1, 5},
		KeyValue{2, 0},
	}

	q := From(outer).GroupJoin(
		From(inner),
		func(i interface{}) interface{} { return i },
		func(i interface{}) interface{} { return i.(int) % 2 },
		func(outer interface{}, inners []interface{}) interface{} {
			return KeyValue{outer, len(inners)}
		})

	if !validateQuery(q, want) {
		t.Errorf("From().GroupJoin()=%v expected %v", toSlice(q), want)
	}
}
