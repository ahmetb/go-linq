package linq

import (
	"testing"
)

func TestDefaultIfEmpty(t *testing.T) {
	defaultValue := 0
	tests := []struct {
		input []any
		want  []any
	}{
		{[]any{}, []any{defaultValue}},
		{[]any{1, 2, 3, 4, 5}, []any{1, 2, 3, 4, 5}},
	}

	for _, test := range tests {
		q := From(test.input).DefaultIfEmpty(defaultValue)

		if !testQueryIteration(q, test.want) {
			t.Errorf("From(%v).DefaultIfEmpty(%v)=%v expected %v", test.input, defaultValue, toSlice(q), test.want)
		}
	}

}
