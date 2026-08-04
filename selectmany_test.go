package linq

import (
	"strconv"
	"testing"
)

func TestSelectMany(t *testing.T) {
	tests := []struct {
		input    any
		selector func(any) legacyQuery
		output   []any
	}{
		{[][]int{{1, 2, 3}, {4, 5, 6, 7}}, func(i any) legacyQuery {
			return legacyFrom(i)
		}, []any{1, 2, 3, 4, 5, 6, 7}},
		{[]string{"str", "ing"}, func(i any) legacyQuery {
			return legacyFromString(i.(string))
		}, []any{'s', 't', 'r', 'i', 'n', 'g'}},
	}

	for _, test := range tests {
		if q := legacyFrom(test.input).SelectMany(test.selector); !testQueryIteration(q, test.output) {
			t.Errorf("From(%v).SelectMany()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestSelectManyT_PanicWhenSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "SelectManyT: parameter [selectorFn] has a invalid function signature. Expected: 'func(T)linq.legacyQuery', actual: 'func(int)int'", func() {
		legacyFrom([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).SelectManyT(func(item int) int { return item + 2 })
	})
}

func TestSelectManyIndexed(t *testing.T) {
	tests := []struct {
		input    any
		selector func(int, any) legacyQuery
		output   []any
	}{
		{[][]int{{1, 2, 3}, {4, 5, 6, 7}}, func(i int, x any) legacyQuery {
			if i > 0 {
				return legacyFrom(x.([]int)[1:])
			}
			return legacyFrom(x)
		}, []any{1, 2, 3, 5, 6, 7}},
		{[]string{"str", "ing"}, func(i int, x any) legacyQuery {
			return legacyFromString(x.(string) + strconv.Itoa(i))
		}, []any{'s', 't', 'r', '0', 'i', 'n', 'g', '1'}},
	}

	for _, test := range tests {
		if q := legacyFrom(test.input).SelectManyIndexed(test.selector); !testQueryIteration(q, test.output) {
			t.Errorf("From(%v).SelectManyIndexed()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestSelectManyIndexedT_PanicWhenSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "SelectManyIndexedT: parameter [selectorFn] has a invalid function signature. Expected: 'func(int,T)linq.legacyQuery', actual: 'func(int)int'", func() {
		legacyFrom([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).SelectManyIndexedT(func(item int) int { return item + 2 })
	})
}

func TestSelectManyBy(t *testing.T) {
	tests := []struct {
		input          any
		selector       func(any) legacyQuery
		resultSelector func(any, any) any
		output         []any
	}{
		{[][]int{{1, 2, 3}, {4, 5, 6, 7}}, func(i any) legacyQuery {
			return legacyFrom(i)
		}, func(x any, y any) any {
			return x.(int) + 1
		}, []any{2, 3, 4, 5, 6, 7, 8}},
		{[]string{"str", "ing"}, func(i any) legacyQuery {
			return legacyFromString(i.(string))
		}, func(x any, y any) any {
			return string(x.(rune)) + "_"
		}, []any{"s_", "t_", "r_", "i_", "n_", "g_"}},
	}

	for _, test := range tests {
		if q := legacyFrom(test.input).SelectManyBy(test.selector, test.resultSelector); !testQueryIteration(q, test.output) {
			t.Errorf("From(%v).SelectManyBy()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestSelectManyByT_PanicWhenSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "SelectManyByT: parameter [selectorFn] has a invalid function signature. Expected: 'func(T)linq.legacyQuery', actual: 'func(int)interface {}'", func() {
		legacyFrom([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).SelectManyByT(func(item int) any { return item + 2 }, 2)
	})
}

func TestSelectManyByT_PanicWhenResultSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "SelectManyByT: parameter [resultSelectorFn] has a invalid function signature. Expected: 'func(T,T)T', actual: 'func()'", func() {
		legacyFrom([][]int{{1, 1, 1, 2}, {1, 2, 3, 4, 2}}).SelectManyByT(
			func(item any) legacyQuery { return legacyFrom(item) },
			func() {},
		)
	})
}

func TestSelectManyIndexedBy(t *testing.T) {
	tests := []struct {
		input          any
		selector       func(int, any) legacyQuery
		resultSelector func(any, any) any
		output         []any
	}{
		{[][]int{{1, 2, 3}, {4, 5, 6, 7}}, func(i int, x any) legacyQuery {
			if i == 0 {
				return legacyFrom([]int{10, 20, 30})
			}
			return legacyFrom(x)
		}, func(x any, y any) any {
			return x.(int) + 1
		}, []any{11, 21, 31, 5, 6, 7, 8}},
		{[]string{"st", "ng"}, func(i int, x any) legacyQuery {
			if i == 0 {
				return legacyFromString(x.(string) + "r")
			}
			return legacyFromString("i" + x.(string))
		}, func(x any, y any) any {
			return string(x.(rune)) + "_"
		}, []any{"s_", "t_", "r_", "i_", "n_", "g_"}},
	}

	for _, test := range tests {
		if q := legacyFrom(test.input).SelectManyByIndexed(test.selector, test.resultSelector); !testQueryIteration(q, test.output) {
			t.Errorf("From(%v).SelectManyIndexedBy()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestSelectManyIndexedByT_PanicWhenSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "SelectManyByIndexedT: parameter [selectorFn] has a invalid function signature. Expected: 'func(int,T)linq.legacyQuery', actual: 'func(int)interface {}'", func() {
		legacyFrom([][]int{{1, 1, 1, 2}, {1, 2, 3, 4, 2}}).SelectManyByIndexedT(
			func(item int) any { return item + 2 },
			2,
		)
	})
}

func TestSelectManyIndexedByT_PanicWhenResultSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "SelectManyByIndexedT: parameter [resultSelectorFn] has a invalid function signature. Expected: 'func(T,T)T', actual: 'func()'", func() {
		legacyFrom([][]int{{1, 1, 1, 2}, {1, 2, 3, 4, 2}}).SelectManyByIndexedT(
			func(index int, item any) legacyQuery { return legacyFrom(item) },
			func() {},
		)
	})
}
