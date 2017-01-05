package linq

import "testing"

func TestTake(t *testing.T) {
	tests := []struct {
		input  interface{}
		output []interface{}
	}{
		{[]int{1, 2, 2, 3, 1}, []interface{}{1, 2, 2}},
		{[9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, []interface{}{1, 1, 1}},
		{"sstr", []interface{}{'s', 's', 't'}},
	}

	for _, test := range tests {
		if q := From(test.input).Take(3); !validateQuery(q, test.output) {
			t.Errorf("From(%v).Take(3)=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestTakeWhile(t *testing.T) {
	tests := []struct {
		input     interface{}
		predicate func(interface{}) bool
		output    []interface{}
	}{
		{[]int{1, 1, 1, 2, 1, 2}, func(i interface{}) bool {
			return i.(int) < 3
		}, []interface{}{1, 1, 1, 2, 1, 2}},
		{[9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, func(i interface{}) bool {
			return i.(int) < 3
		}, []interface{}{1, 1, 1, 2, 1, 2}},
		{"sstr", func(i interface{}) bool {
			return i.(rune) == 's'
		}, []interface{}{'s', 's'}},
	}

	for _, test := range tests {
		if q := From(test.input).TakeWhile(test.predicate); !validateQuery(q, test.output) {
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
		input     interface{}
		predicate func(int, interface{}) bool
		output    []interface{}
	}{
		{[]int{1, 1, 1, 2}, func(i int, x interface{}) bool {
			return x.(int) < 2 || i < 5
		}, []interface{}{1, 1, 1, 2}},
		{[9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, func(i int, x interface{}) bool {
			return x.(int) < 2 || i < 5
		}, []interface{}{1, 1, 1, 2, 1}},
		{"sstr", func(i int, x interface{}) bool {
			return x.(rune) == 's' && i < 1
		}, []interface{}{'s'}},
	}

	for _, test := range tests {
		if q := From(test.input).TakeWhileIndexed(test.predicate); !validateQuery(q, test.output) {
			t.Errorf("From(%v).TakeWhileIndexed()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestTakeWhileIndexedT_PanicWhenPredicateFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "TakeWhileIndexedT: parameter [predicateFn] has a invalid function signature. Expected: 'func(int,T)bool', actual: 'func(int)int'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).TakeWhileIndexedT(func(item int) int { return item + 2 })
	})
}
