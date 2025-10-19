package linq

import "testing"

func TestDistinct(t *testing.T) {
	tests := []struct {
		input  any
		output []any
	}{
		{[]int{1, 2, 2, 3, 1}, []any{1, 2, 3}},
		{[9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, []any{1, 2, 3, 4}},
		{"sstr", []any{'s', 't', 'r'}},
	}

	for _, test := range tests {
		if q := From(test.input).Distinct(); !testQueryIteration(q, test.output) {
			t.Errorf("From(%v).Distinct()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestDistinctForOrderedQuery(t *testing.T) {
	tests := []struct {
		input  any
		output []any
	}{
		{[]int{1, 2, 2, 3, 1}, []any{1, 2, 3}},
		{[9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, []any{1, 2, 3, 4}},
		{"sstr", []any{'r', 's', 't'}},
	}

	for _, test := range tests {
		if q := From(test.input).OrderBy(func(i any) any {
			return i
		}).Distinct(); !testQueryIteration(q.Query, test.output) {
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
	want := []any{user{1, "Foo"}, user{2, "Bar"}}

	if q := From(users).DistinctBy(func(u any) any {
		return u.(user).name
	}); !testQueryIteration(q, want) {
		t.Errorf("From(%v).DistinctBy()=%v expected %v", users, toSlice(q), want)
	}
}

func TestDistinctByT_PanicWhenSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "DistinctByT: parameter [selectorFn] has a invalid function signature. Expected: 'func(T)T', actual: 'func(string,string)bool'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).DistinctByT(func(indice, item string) bool { return item == "2" })
	})
}
