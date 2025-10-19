package linq

import (
	"strconv"
	"testing"
)

func TestSelect(t *testing.T) {
	tests := []struct {
		input    any
		selector func(any) any
		output   []any
	}{
		{[]int{1, 2, 3}, func(i any) any {
			return i.(int) * 2
		}, []any{2, 4, 6}},
		{"str", func(i any) any {
			return string(i.(rune)) + "1"
		}, []any{"s1", "t1", "r1"}},
	}

	for _, test := range tests {
		if q := From(test.input).Select(test.selector); !testQueryIteration(q, test.output) {
			t.Errorf("From(%v).Select()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestSelectT_PanicWhenSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "SelectT: parameter [selectorFn] has a invalid function signature. Expected: 'func(T)T', actual: 'func(int,int)int'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).SelectT(func(item, idx int) int { return item + 2 })
	})
}

func TestSelectIndexed(t *testing.T) {
	tests := []struct {
		input    any
		selector func(int, any) any
		output   []any
	}{
		{[]int{1, 2, 3}, func(i int, x any) any {
			return x.(int) * i
		}, []any{0, 2, 6}},
		{"str", func(i int, x any) any {
			return string(x.(rune)) + strconv.Itoa(i)
		}, []any{"s0", "t1", "r2"}},
	}

	for _, test := range tests {
		if q := From(test.input).SelectIndexed(test.selector); !testQueryIteration(q, test.output) {
			t.Errorf("From(%v).SelectIndexed()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestSelectIndexedT_PanicWhenSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "SelectIndexedT: parameter [selectorFn] has a invalid function signature. Expected: 'func(int,T)T', actual: 'func(string,int)int'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).SelectIndexedT(func(index string, item int) int { return item + 2 })
	})
}
