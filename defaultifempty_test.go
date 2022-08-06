package linq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultIfEmpty(t *testing.T) {
	defaultValue := 0
	tests := []struct {
		input []interface{}
		want  []interface{}
	}{
		{[]interface{}{}, []interface{}{defaultValue}},
		{[]interface{}{1, 2, 3, 4, 5}, []interface{}{1, 2, 3, 4, 5}},
	}

	for _, test := range tests {
		q := From(test.input).DefaultIfEmpty(defaultValue)

		if !validateQuery(q, test.want) {
			t.Errorf("From(%v).DefaultIfEmpty(%v)=%v expected %v", test.input, defaultValue, toSlice(q), test.want)
		}
	}

}

func TestDefaultIfEmptyG(t *testing.T) {
	defaultValue := 0
	tests := []struct {
		input []int
		want  []int
	}{
		{[]int{}, []int{defaultValue}},
		{[]int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
	}

	for _, test := range tests {
		actual := FromSliceG(test.input).DefaultIfEmpty(defaultValue).ToSlice()

		assert.Equal(t, test.want, actual)
	}

}
