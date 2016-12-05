package linq

import "testing"

func TestTakeWhileT(t *testing.T) {
	tests := []struct {
		input     interface{}
		predicate interface{}
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
		{"sstr", func(i rune) bool {
			return i == 's'
		}, []interface{}{'s', 's'}},
	}

	for _, test := range tests {
		if q := From(test.input).TakeWhileT(test.predicate); !validateQuery(q, test.output) {
			t.Errorf("From(%v).TakeWhile()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestTakeWhileT_PanicWhenFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).TakeWhileT(func(item int) int { return item + 2 })
}

func TestTakeWhileIndexedT(t *testing.T) {
	tests := []struct {
		input     interface{}
		predicate interface{}
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
		{"sstr", func(i int, x rune) bool {
			return x == 's' && i < 1
		}, []interface{}{'s'}},
	}

	for _, test := range tests {
		if q := From(test.input).TakeWhileIndexedT(test.predicate); !validateQuery(q, test.output) {
			t.Errorf("From(%v).TakeWhileIndexed()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestTakeWhileIndexedT_PanicWhenFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).TakeWhileIndexedT(func(item int) int { return item + 2 })
}
