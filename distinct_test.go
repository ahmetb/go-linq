package linq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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

func TestDistinctG(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3}, FromSliceG([]int{1, 2, 2, 3, 1}).Distinct().ToSlice())
	assert.Equal(t, []int{1, 2, 3, 4}, FromSliceG([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).Distinct().ToSlice())
	assert.Equal(t, []rune{'s', 't', 'r'}, FromStringG("sstr").Distinct().ToSlice())
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

func TestDistinctForOrderedQueryG(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3}, FromSliceG([]int{1, 2, 2, 3, 1}).Expend(To2[int, int]()).(Expended[int, int]).OrderBy(func(i int) int {
		return i
	}).Distinct().ToSlice())
	assert.Equal(t, []int{1, 2, 3, 4}, FromSliceG([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).Expend(To2[int, int]()).(Expended[int, int]).OrderBy(func(i int) int {
		return i
	}).Distinct().ToSlice())
	assert.Equal(t, []rune{'r', 's', 't'}, FromStringG("sstr").Expend(To2[rune, rune]()).(Expended[rune, rune]).OrderBy(func(i rune) rune {
		return i
	}).Distinct().ToSlice())
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

func TestDistinctByG(t *testing.T) {
	type user struct {
		id   int
		name string
	}

	users := []user{{1, "Foo"}, {2, "Bar"}, {3, "Foo"}}
	want := []user{user{1, "Foo"}, user{2, "Bar"}}

	assert.Equal(t, want, FromSliceG(users).Expend(To2[user, string]()).(Expended[user, string]).DistinctBy(func(u user) string {
		return u.name
	}).ToSlice())
}

func TestDistinctByT_PanicWhenSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "DistinctByT: parameter [selectorFn] has a invalid function signature. Expected: 'func(T)T', actual: 'func(string,string)bool'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).DistinctByT(func(indice, item string) bool { return item == "2" })
	})
}
