package linq

import (
	"strconv"
	"testing"
)

func TestSelectManyT(t *testing.T) {
	tests := []struct {
		input    interface{}
		selector interface{}
		output   []interface{}
	}{
		{[][]int{{1, 2, 3}, {4, 5, 6, 7}}, func(i interface{}) Query {
			return From(i)
		}, []interface{}{1, 2, 3, 4, 5, 6, 7}},
		{[][]int{{1, 2, 3}, {4, 5, 6, 7}}, func(i []int) Query {
			return From(i)
		}, []interface{}{1, 2, 3, 4, 5, 6, 7}},
		{[]string{"str", "ing"}, func(i interface{}) Query {
			return FromString(i.(string))
		}, []interface{}{'s', 't', 'r', 'i', 'n', 'g'}},
	}

	for _, test := range tests {
		if q := From(test.input).SelectManyT(test.selector); !validateQuery(q, test.output) {
			t.Errorf("From(%v).SelectManyT()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestSelectManyT_PanicWhenFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}
	}()

	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).SelectManyT(func(item int) int { return item + 2 })
}

func TestSelectManyIndexedT(t *testing.T) {
	tests := []struct {
		input    interface{}
		selector interface{}
		output   []interface{}
	}{
		{[][]int{{1, 2, 3}, {4, 5, 6, 7}}, func(i int, x interface{}) Query {
			if i > 0 {
				return From(x.([]int)[1:])
			}
			return From(x)
		}, []interface{}{1, 2, 3, 5, 6, 7}},
		{[][]int{{1, 2, 3}, {4, 5, 6, 7}}, func(i int, x []int) Query {
			if i > 0 {
				return From(x[1:])
			}
			return From(x)
		}, []interface{}{1, 2, 3, 5, 6, 7}},
		{[]string{"str", "ing"}, func(i int, x string) Query {
			return FromString(x + strconv.Itoa(i))
		}, []interface{}{'s', 't', 'r', '0', 'i', 'n', 'g', '1'}},
	}

	for _, test := range tests {
		if q := From(test.input).SelectManyIndexedT(test.selector); !validateQuery(q, test.output) {
			t.Errorf("From(%v).SelectManyIndexedT()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestSelectManyIndexedT_PanicWhenFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}
	}()

	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).SelectManyIndexedT(func(item int) int { return item + 2 })
}

func TestSelectManyByT(t *testing.T) {
	tests := []struct {
		input          interface{}
		selector       interface{}
		resultSelector interface{}
		output         []interface{}
	}{
		{[][]int{{1, 2, 3}, {4, 5, 6, 7}}, func(i []int) Query {
			return From(i)
		}, func(x int, y []int) int {
			return x + 1
		}, []interface{}{2, 3, 4, 5, 6, 7, 8}},
		{[]string{"str", "ing"}, func(i string) Query {
			return FromString(i)
		}, func(x rune, y string) string {
			return string(x) + "_"
		}, []interface{}{"s_", "t_", "r_", "i_", "n_", "g_"}},
	}

	for _, test := range tests {
		if q := From(test.input).SelectManyByT(test.selector, test.resultSelector); !validateQuery(q, test.output) {
			t.Errorf("From(%v).SelectManyBy()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestSelectManyByT_PanicWhenSelectorFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}
	}()

	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).SelectManyByT(func(item int) interface{} { return item + 2 }, 2)
}

func TestSelectManyByT_PanicWhenResultSelectorFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}
	}()

	From([][]int{{1, 1, 1, 2}, {1, 2, 3, 4, 2}}).SelectManyByT(
		func(item interface{}) Query { return From(item) },
		func() {},
	)
}

func TestSelectManyIndexedByT(t *testing.T) {
	tests := []struct {
		input          interface{}
		selector       interface{}
		resultSelector interface{}
		output         []interface{}
	}{
		{[][]int{{1, 2, 3}, {4, 5, 6, 7}}, func(i int, x []int) Query {
			if i == 0 {
				return From([]int{10, 20, 30})
			}
			return From(x)
		}, func(x int, y []int) int {
			return x + 1
		}, []interface{}{11, 21, 31, 5, 6, 7, 8}},
		{[]string{"st", "ng"}, func(i int, x string) Query {
			if i == 0 {
				return FromString(x + "r")
			}
			return FromString("i" + x)
		}, func(x rune, y string) string {
			return string(x) + "_"
		}, []interface{}{"s_", "t_", "r_", "i_", "n_", "g_"}},
	}

	for _, test := range tests {
		if q := From(test.input).SelectManyByIndexedT(test.selector, test.resultSelector); !validateQuery(q, test.output) {
			t.Errorf("From(%v).SelectManyIndexedByT()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestSelectManyIndexedByT_PanicWhenSelectorFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}
	}()

	From([][]int{{1, 1, 1, 2}, {1, 2, 3, 4, 2}}).SelectManyByIndexedT(
		func(item int) interface{} { return item + 2 },
		2,
	)
}

func TestSelectManyIndexedByT_PanicWhenResultSelectorFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}
	}()

	From([][]int{{1, 1, 1, 2}, {1, 2, 3, 4, 2}}).SelectManyByIndexedT(
		func(index int, item interface{}) Query { return From(item) },
		func() {},
	)
}
