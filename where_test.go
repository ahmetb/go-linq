package linq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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

func TestWhereG(t *testing.T) {
	inputs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	q := FromSliceG(inputs)
	greaterThanFive := q.Where(func(item int) bool {
		return item > 5
	}).ToSlice()
	expected := []int{6, 7, 8, 9, 10}
	assert.Equal(t, expected, greaterThanFive)
}

func TestWhereT_PanicWhenPredicateFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "WhereT: parameter [predicateFn] has a invalid function signature. Expected: 'func(T)bool', actual: 'func(int)int'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).WhereT(func(item int) int { return item + 2 })
	})
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
		{"abcde", func(i int, _ interface{}) bool {
			return i < 2
		}, []interface{}{'a', 'b'}},
	}

	for _, test := range tests {
		if q := From(test.input).WhereIndexed(test.predicate); !validateQuery(q, test.output) {
			t.Errorf("From(%v).WhereIndexed()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestWhereIndexedG(t *testing.T) {
	assert.Equal(t, []int{2, 3, 2}, FromSliceG([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).WhereIndexed(func(i, x int) bool {
		return x < 4 && i > 4
	}).ToSlice())
	assert.Equal(t, []rune{'s', 't', 'r'}, FromStringG("sstr").WhereIndexed(func(i int, x rune) bool {
		return x != 's' || i == 1
	}).ToSlice())
	assert.Equal(t, []rune{'a', 'b'}, FromStringG("abcde").WhereIndexed(func(i int, x rune) bool {
		return i < 2
	}).ToSlice())
}

func TestWhereIndexedT_PanicWhenPredicateFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "WhereIndexedT: parameter [predicateFn] has a invalid function signature. Expected: 'func(int,T)bool', actual: 'func(string)'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).WhereIndexedT(func(item string) {})
	})
}
