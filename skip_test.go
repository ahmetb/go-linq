package linq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSkip(t *testing.T) {
	tests := []struct {
		input  interface{}
		output []interface{}
	}{
		{[]int{1, 2}, []interface{}{}},
		{[]int{1, 2, 2, 3, 1}, []interface{}{3, 1}},
		{[9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, []interface{}{2, 1, 2, 3, 4, 2}},
		{"sstr", []interface{}{'r'}},
	}

	for _, test := range tests {
		if q := From(test.input).Skip(3); !validateQuery(q, test.output) {
			t.Errorf("From(%v).Skip(3)=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestSkipG(t *testing.T) {
	assert.Empty(t, FromSliceG([]int{1, 2}).Skip(3).ToSlice())
	assert.Equal(t, []int{3, 1}, FromSliceG([]int{1, 2, 2, 3, 1}).Skip(3).ToSlice())
	assert.Equal(t, []int{2, 1, 2, 3, 4, 2}, FromSliceG([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).Skip(3).ToSlice())
	assert.Equal(t, []rune{'r'}, FromStringG("sstr").Skip(3).ToSlice())
}

func TestSkipWhile(t *testing.T) {
	tests := []struct {
		input     interface{}
		predicate func(interface{}) bool
		output    []interface{}
	}{
		{[]int{1, 2}, func(i interface{}) bool {
			return i.(int) < 3
		}, []interface{}{}},
		{[]int{4, 1, 2}, func(i interface{}) bool {
			return i.(int) < 3
		}, []interface{}{4, 1, 2}},
		{[9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, func(i interface{}) bool {
			return i.(int) < 3
		}, []interface{}{3, 4, 2}},
		{"sstr", func(i interface{}) bool {
			return i.(rune) == 's'
		}, []interface{}{'t', 'r'}},
	}

	for _, test := range tests {
		if q := From(test.input).SkipWhile(test.predicate); !validateQuery(q, test.output) {
			t.Errorf("From(%v).SkipWhile()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestSkipWhileG(t *testing.T) {
	assert.Empty(t, FromSliceG([]int{1, 2}).SkipWhile(func(i int) bool {
		return i < 3
	}).ToSlice())
	assert.Equal(t, []int{4, 1, 2}, FromSliceG([]int{4, 1, 2}).SkipWhile(func(i int) bool {
		return i < 3
	}).ToSlice())
	assert.Equal(t, []int{3, 4, 2}, FromSliceG([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).SkipWhile(func(i int) bool {
		return i < 3
	}).ToSlice())
	assert.Equal(t, []rune{'t', 'r'}, FromStringG("sstr").SkipWhile(func(i rune) bool {
		return i == 's'
	}).ToSlice())
}

func TestSkipWhileT_PanicWhenPredicateFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "SkipWhileT: parameter [predicateFn] has a invalid function signature. Expected: 'func(T)bool', actual: 'func(int,int)bool'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).SkipWhileT(func(item int, x int) bool { return item == 1 })
	})
}

func TestSkipWhileIndexed(t *testing.T) {
	tests := []struct {
		input     interface{}
		predicate func(int, interface{}) bool
		output    []interface{}
	}{
		{[]int{1, 2}, func(i int, x interface{}) bool {
			return x.(int) < 3
		}, []interface{}{}},
		{[]int{4, 1, 2}, func(i int, x interface{}) bool {
			return x.(int) < 3
		}, []interface{}{4, 1, 2}},
		{[9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}, func(i int, x interface{}) bool {
			return x.(int) < 2 || i < 5
		}, []interface{}{2, 3, 4, 2}},
		{"sstr", func(i int, x interface{}) bool {
			return x.(rune) == 's' && i < 1
		}, []interface{}{'s', 't', 'r'}},
	}

	for _, test := range tests {
		if q := From(test.input).SkipWhileIndexed(test.predicate); !validateQuery(q, test.output) {
			t.Errorf("From(%v).SkipWhileIndexed()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestSkipWhileIndexedG(t *testing.T) {
	assert.Empty(t, FromSliceG([]int{1, 2}).SkipWhileIndexed(func(i int, x int) bool {
		return x < 3
	}).ToSlice())
	assert.Equal(t, []int{4, 1, 2}, FromSliceG([]int{4, 1, 2}).SkipWhileIndexed(func(i int, x int) bool {
		return x < 3
	}).ToSlice())
	assert.Equal(t, []int{2, 3, 4, 2}, FromSliceG([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).SkipWhileIndexed(func(i int, x int) bool {
		return x < 2 || i < 5
	}).ToSlice())
	assert.Equal(t, []rune{'s', 't', 'r'}, FromStringG("sstr").SkipWhileIndexed(func(i int, x rune) bool {
		return x == 's' && i < 1
	}).ToSlice())
}

func TestSkipWhileIndexedT_PanicWhenPredicateFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "SkipWhileIndexedT: parameter [predicateFn] has a invalid function signature. Expected: 'func(int,T)bool', actual: 'func(int,int,int)bool'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).SkipWhileIndexedT(func(item int, x int, y int) bool { return item == 1 })
	})
}
