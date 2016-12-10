package linq

import "testing"

func TestWhere(t *testing.T) {
	tests := []struct {
		input     interface{}
		predicate func(interface{}) bool
		output    []interface{}
	}{
		{[9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, func(i interface{}) bool {
			return i.(int) >= 3
		}, []interface{}{3, 4}},
		{"sstr", func(i interface{}) bool {
			return i.(rune) != 's'
		}, []interface{}{'t', 'r'}},
	}

	for _, test := range tests {
		if q := From(test.input).Where(test.predicate); !validateQuery(q, test.output) {
			t.Errorf("From(%v).Where()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestWhereT(t *testing.T) {
	tests := []struct {
		input     interface{}
		predicate interface{}
		output    []interface{}
	}{
		{[9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, func(i interface{}) bool {
			return i.(int) >= 3
		}, []interface{}{3, 4}},
		{"sstr", func(i interface{}) bool {
			return i.(rune) != 's'
		}, []interface{}{'t', 'r'}},
		{"sstr", func(i rune) bool {
			return i != 's'
		}, []interface{}{'t', 'r'}},
	}

	for _, test := range tests {
		if q := From(test.input).WhereT(test.predicate); !validateQuery(q, test.output) {
			t.Errorf("From(%v).WhereT()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestWhereT_PanicWhenFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).WhereT(func(item int) int { return item + 2 })
}

func TestWhereIndexed(t *testing.T) {
	tests := []struct {
		input     interface{}
		predicate func(int, interface{}) bool
		output    []interface{}
	}{
		{[9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, func(i int, x interface{}) bool {
			return x.(int) < 4 && i > 4
		}, []interface{}{2, 3, 2}},
		{"sstr", func(i int, x interface{}) bool {
			return x.(rune) != 's' || i == 1
		}, []interface{}{'s', 't', 'r'}},
	}

	for _, test := range tests {
		if q := From(test.input).WhereIndexed(test.predicate); !validateQuery(q, test.output) {
			t.Errorf("From(%v).WhereIndexed()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestWhereIndexedT(t *testing.T) {
	tests := []struct {
		input     interface{}
		predicate interface{}
		output    []interface{}
	}{
		{[9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, func(i int, x interface{}) bool {
			return x.(int) < 4 && i > 4
		}, []interface{}{2, 3, 2}},
		{"sstr", func(i int, x interface{}) bool {
			return x.(rune) != 's' || i == 1
		}, []interface{}{'s', 't', 'r'}},
		{"sstr", func(i int, x rune) bool {
			return x != 's' || i == 1
		}, []interface{}{'s', 't', 'r'}},
	}

	for _, test := range tests {
		if q := From(test.input).WhereIndexedT(test.predicate); !validateQuery(q, test.output) {
			t.Errorf("From(%v).WhereIndexed()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestWhereIndexedT_PanicWhenFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic")
		}
	}()

	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).WhereIndexedT(func(item string) {})
}

func TestWhereIndexedT_PanicWhenParameterTypesDoesntMatch(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic")
		}
	}()
	var result []int
	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).WhereIndexedT(func(index int, item string) bool { return true }).ToSlice(&result)
}
