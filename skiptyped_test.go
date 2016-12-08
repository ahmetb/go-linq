package linq

import "testing"

func TestSkipWhileT(t *testing.T) {
	tests := []struct {
		input     interface{}
		predicate interface{}
		output    []interface{}
	}{
		{[]int{1, 2}, func(i interface{}) bool {
			return i.(int) < 3
		}, []interface{}{}},
		{[]int{4, 1, 2}, func(i interface{}) bool {
			return i.(int) < 3
		}, []interface{}{4, 1, 2}},
		{[9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, func(i int) bool {
			return i < 3
		}, []interface{}{3, 4, 2}},
		{"sstr", func(i interface{}) bool {
			return i.(rune) == 's'
		}, []interface{}{'t', 'r'}},
	}

	for _, test := range tests {
		if q := From(test.input).SkipWhileT(test.predicate); !validateQuery(q, test.output) {
			t.Errorf("From(%v).SkipWhile()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestSkipWhileT_PanicWhenFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).SkipWhileT(func(item int, x int) bool { return item == 1 })
}

func TestSkipWhileIndexedT(t *testing.T) {
	tests := []struct {
		input     interface{}
		predicate interface{}
		output    []interface{}
	}{
		{[]int{1, 2}, func(i int, x interface{}) bool {
			return x.(int) < 3
		}, []interface{}{}},
		{[]int{4, 1, 2}, func(i int, x interface{}) bool {
			return x.(int) < 3
		}, []interface{}{4, 1, 2}},
		{[9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, func(i int, x int) bool {
			return x < 2 || i < 5
		}, []interface{}{2, 3, 4, 2}},
		{"sstr", func(i int, x interface{}) bool {
			return x.(rune) == 's' && i < 1
		}, []interface{}{'s', 't', 'r'}},
	}

	for _, test := range tests {
		if q := From(test.input).SkipWhileIndexedT(test.predicate); !validateQuery(q, test.output) {
			t.Errorf("From(%v).SkipWhileIndexed()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestSkipWhileIndexedT_PanicWhenFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).SkipWhileIndexedT(func(item int, x int, y int) bool { return item == 1 })
}
