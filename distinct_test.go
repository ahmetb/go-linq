package linq

import "testing"

func TestDistinct(t *testing.T) {
	arr := [9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}

	tests := []struct {
		input  Query[int]
		output []int
	}{
		{FromSlice([]int{1, 2, 2, 3, 1}), []int{1, 2, 3}},
		{FromSlice(arr[:]), []int{1, 2, 3, 4}},
	}

	for _, test := range tests {
		if q := test.input.Distinct(); !testQueryIteration(q, test.output) {
			t.Errorf("Distinct()=%v expected %v", toSlice(q), test.output)
		}
	}

	want := []rune{'s', 't', 'r'}
	if q := FromString("sstr").Distinct(); !testQueryIteration(q, want) {
		t.Errorf("FromString(sstr).Distinct()=%v expected %v", toSlice(q), want)
	}
}

func TestDistinctForOrderedQuery(t *testing.T) {
	arr := [9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}

	tests := []struct {
		input  Query[int]
		output []int
	}{
		{FromSlice([]int{1, 2, 2, 3, 1}), []int{1, 2, 3}},
		{FromSlice(arr[:]), []int{1, 2, 3, 4}},
	}

	for _, test := range tests {
		if q := test.input.OrderBy(func(i int) int {
			return i
		}).Distinct(); !testQueryIteration(q.Query, test.output) {
			t.Errorf("Distinct()=%v expected %v", toSlice(q.Query), test.output)
		}
	}

	want := []rune{'r', 's', 't'}
	if q := FromString("sstr").OrderBy(func(r rune) rune {
		return r
	}).Distinct(); !testQueryIteration(q.Query, want) {
		t.Errorf("FromString(sstr).OrderBy().Distinct()=%v expected %v", toSlice(q.Query), want)
	}
}

func TestDistinctBy(t *testing.T) {
	type user struct {
		id   int
		name string
	}

	users := []user{{1, "Foo"}, {2, "Bar"}, {3, "Foo"}}
	want := []user{{1, "Foo"}, {2, "Bar"}}

	if q := FromSlice(users).DistinctBy(func(u user) string {
		return u.name
	}); !testQueryIteration(q, want) {
		t.Errorf("DistinctBy()=%v expected %v", toSlice(q), want)
	}
}
