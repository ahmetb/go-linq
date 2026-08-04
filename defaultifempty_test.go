package linq

import (
	"testing"
)

func TestDefaultIfEmpty(t *testing.T) {
	defaultValue := 0
	tests := []struct {
		input []int
		want  []int
	}{
		{[]int{}, []int{defaultValue}},
		{[]int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
	}

	for _, test := range tests {
		q := FromSlice(test.input).DefaultIfEmpty(defaultValue)

		if !testQueryIteration(q, test.want) {
			t.Errorf("FromSlice(%v).DefaultIfEmpty(%v)=%v expected %v", test.input, defaultValue, toSlice(q), test.want)
		}
	}
}
