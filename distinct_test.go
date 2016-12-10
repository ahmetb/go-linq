package linq

import "testing"

func TestDistinct(t *testing.T) {
	tests := []struct {
		input  interface{}
		output []interface{}
	}{
		{[]int{1, 2, 2, 3, 1}, []interface{}{1, 2, 3}},
		{[9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, []interface{}{1, 2, 3, 4}},
		{"sstr", []interface{}{'s', 't', 'r'}},
	}

	for _, test := range tests {
		if q := From(test.input).Distinct(); !validateQuery(q, test.output) {
			t.Errorf("From(%v).Distinct()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestDistinctForOrderedQuery(t *testing.T) {
	tests := []struct {
		input  interface{}
		output []interface{}
	}{
		{[]int{1, 2, 2, 3, 1}, []interface{}{1, 2, 3}},
		{[9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, []interface{}{1, 2, 3, 4}},
		{"sstr", []interface{}{'r', 's', 't'}},
	}

	for _, test := range tests {
		if q := From(test.input).OrderBy(func(i interface{}) interface{} {
			return i
		}).Distinct(); !validateQuery(q.Query, test.output) {
			t.Errorf("From(%v).Distinct()=%v expected %v", test.input, toSlice(q.Query), test.output)
		}
	}
}

func TestDistinctBy(t *testing.T) {
	type user struct {
		id   int
		name string
	}

	users := []user{{1, "Foo"}, {2, "Bar"}, {3, "Foo"}}
	want := []interface{}{user{1, "Foo"}, user{2, "Bar"}}

	if q := From(users).DistinctBy(func(u interface{}) interface{} {
		return u.(user).name
	}); !validateQuery(q, want) {
		t.Errorf("From(%v).DistinctBy()=%v expected %v", users, toSlice(q), want)
	}
}

func TestDistinctByT(t *testing.T) {
	type user struct {
		id   int
		name string
	}

	tests := []struct {
		input    interface{}
		selector interface{}
		output   []interface{}
	}{
		{[]user{{1, "Foo"}, {2, "Bar"}, {3, "Foo"}}, func(u interface{}) interface{} {
			return u.(user).name
		}, []interface{}{user{1, "Foo"}, user{2, "Bar"}}},
		{[]user{{1, "Foo"}, {2, "Bar"}, {3, "Foo"}}, func(u user) string {
			return u.name
		}, []interface{}{user{1, "Foo"}, user{2, "Bar"}}},
	}

	for _, test := range tests {
		if q := From(test.input).DistinctByT(test.selector); !validateQuery(q, test.output) {
			t.Errorf("From(%v).DistinctByT()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestDistinctBy_PanicWhenFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).DistinctByT(func(indice, item string) bool { return item == "2" })
}
