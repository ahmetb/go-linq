package linq

import "testing"

func TestReverse(t *testing.T) {
	tests := []struct {
		input interface{}
		want  []interface{}
	}{
		{[]int{1, 2, 3}, []interface{}{3, 2, 1}},
	}

	for _, test := range tests {
		if q := From(test.input).Reverse(); !validateQuery(q, test.want) {
			t.Errorf("From(%v).Reverse()=%v expected %v", test.input, toSlice(q), test.want)
		}
	}
}
