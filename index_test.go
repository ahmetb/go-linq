package linq

import (
	"testing"
)

func TestIndexOf(t *testing.T) {
	tests := []struct {
		input     interface{}
		predicate func(interface{}) bool
		expected  int
	}{
		{
			input: [9]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			predicate: func(i interface{}) bool {
				return i.(int) == 3
			},
			expected: 2,
		},
		{
			input: "sstr",
			predicate: func(i interface{}) bool {
				return i.(rune) == 'r'
			},
			expected: 3,
		},
		{
			input: "gadsgsadgsda",
			predicate: func(i interface{}) bool {
				return i.(rune) == 'z'
			},
			expected: -1,
		},
	}

	for _, test := range tests {
		index := From(test.input).IndexOf(test.predicate)
		if index != test.expected {
			t.Errorf("From(%v).IndexOf() expected %v received %v", test.input, test.expected, index)
		}

		index = From(test.input).IndexOfT(test.predicate)
		if index != test.expected {
			t.Errorf("From(%v).IndexOfT() expected %v received %v", test.input, test.expected, index)
		}
	}
}

func TestIndexOfT_PanicWhenPredicateFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "IndexOfT: parameter [predicateFn] has a invalid function signature. Expected: 'func(T)bool', actual: 'func(int)int'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).IndexOfT(func(item int) int { return item + 2 })
	})
}
