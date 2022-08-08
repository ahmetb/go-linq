package linq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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

func TestReverseG(t *testing.T) {
	input := []int{1, 2, 3}
	assert.Equal(t, []int{3, 2, 1}, FromSliceG(input).Reverse().ToSlice())
}
