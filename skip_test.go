package linq

import "testing"

func TestSkip(t *testing.T) {
	tests := []struct {
		input  any
		output []any
	}{
		{[]int{1, 2}, []any{}},
		{[]int{1, 2, 2, 3, 1}, []any{3, 1}},
		{[9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, []any{2, 1, 2, 3, 4, 2}},
		{"sstr", []any{'r'}},
	}

	for _, test := range tests {
		if q := From(test.input).Skip(3); !validateQuery(q, test.output) {
			t.Errorf("From(%v).Skip(3)=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestSkipWhile(t *testing.T) {
	tests := []struct {
		input     any
		predicate func(any) bool
		output    []any
	}{
		{[]int{1, 2}, func(i any) bool {
			return i.(int) < 3
		}, []any{}},
		{[]int{4, 1, 2}, func(i any) bool {
			return i.(int) < 3
		}, []any{4, 1, 2}},
		{[9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, func(i any) bool {
			return i.(int) < 3
		}, []any{3, 4, 2}},
		{"sstr", func(i any) bool {
			return i.(rune) == 's'
		}, []any{'t', 'r'}},
	}

	for _, test := range tests {
		if q := From(test.input).SkipWhile(test.predicate); !validateQuery(q, test.output) {
			t.Errorf("From(%v).SkipWhile()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestSkipWhileT_PanicWhenPredicateFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "SkipWhileT: parameter [predicateFn] has a invalid function signature. Expected: 'func(T)bool', actual: 'func(int,int)bool'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).SkipWhileT(func(item int, x int) bool { return item == 1 })
	})
}

func TestSkipWhileIndexed(t *testing.T) {
	tests := []struct {
		input     any
		predicate func(int, any) bool
		output    []any
	}{
		{[]int{1, 2}, func(i int, x any) bool {
			return x.(int) < 3
		}, []any{}},
		{[]int{4, 1, 2}, func(i int, x any) bool {
			return x.(int) < 3
		}, []any{4, 1, 2}},
		{[9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, func(i int, x any) bool {
			return x.(int) < 2 || i < 5
		}, []any{2, 3, 4, 2}},
		{"sstr", func(i int, x any) bool {
			return x.(rune) == 's' && i < 1
		}, []any{'s', 't', 'r'}},
	}

	for _, test := range tests {
		if q := From(test.input).SkipWhileIndexed(test.predicate); !validateQuery(q, test.output) {
			t.Errorf("From(%v).SkipWhileIndexed()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestSkipWhileIndexedT_PanicWhenPredicateFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "SkipWhileIndexedT: parameter [predicateFn] has a invalid function signature. Expected: 'func(int,T)bool', actual: 'func(int,int,int)bool'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).SkipWhileIndexedT(func(item int, x int, y int) bool { return item == 1 })
	})
}
