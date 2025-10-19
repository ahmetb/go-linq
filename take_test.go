package linq

import "testing"

func TestTake(t *testing.T) {
	tests := []struct {
		input  any
		output []any
	}{
		{[]int{1, 2, 2, 3, 1}, []any{1, 2, 2}},
		{[9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, []any{1, 1, 1}},
		{"sstr", []any{'s', 's', 't'}},
	}

	for _, test := range tests {
		if q := From(test.input).Take(3); !testQueryIteration(q, test.output) {
			t.Errorf("From(%v).Take(3)=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestTakeWhile(t *testing.T) {
	tests := []struct {
		input     any
		predicate func(any) bool
		output    []any
	}{
		{[]int{1, 1, 1, 2, 1, 2}, func(i any) bool {
			return i.(int) < 3
		}, []any{1, 1, 1, 2, 1, 2}},
		{[9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, func(i any) bool {
			return i.(int) < 3
		}, []any{1, 1, 1, 2, 1, 2}},
		{"sstr", func(i any) bool {
			return i.(rune) == 's'
		}, []any{'s', 's'}},
	}

	for _, test := range tests {
		if q := From(test.input).TakeWhile(test.predicate); !testQueryIteration(q, test.output) {
			t.Errorf("From(%v).TakeWhile()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestTakeWhileT_PanicWhenPredicateFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "TakeWhileT: parameter [predicateFn] has a invalid function signature. Expected: 'func(T)bool', actual: 'func(int)int'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).TakeWhileT(func(item int) int { return item + 2 })
	})
}

func TestTakeWhileIndexed(t *testing.T) {
	tests := []struct {
		input     any
		predicate func(int, any) bool
		output    []any
	}{
		{[]int{1, 1, 1, 2}, func(i int, x any) bool {
			return x.(int) < 2 || i < 5
		}, []any{1, 1, 1, 2}},
		{[9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, func(i int, x any) bool {
			return x.(int) < 2 || i < 5
		}, []any{1, 1, 1, 2, 1}},
		{"sstr", func(i int, x any) bool {
			return x.(rune) == 's' && i < 1
		}, []any{'s'}},
	}

	for _, test := range tests {
		if q := From(test.input).TakeWhileIndexed(test.predicate); !testQueryIteration(q, test.output) {
			t.Errorf("From(%v).TakeWhileIndexed()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestTakeWhileIndexedT_PanicWhenPredicateFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "TakeWhileIndexedT: parameter [predicateFn] has a invalid function signature. Expected: 'func(int,T)bool', actual: 'func(int)int'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).TakeWhileIndexedT(func(item int) int { return item + 2 })
	})
}
