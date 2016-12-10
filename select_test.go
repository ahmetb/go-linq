package linq

import (
	"strconv"
	"testing"
)

func TestSelect(t *testing.T) {
	tests := []struct {
		input    interface{}
		selector func(interface{}) interface{}
		output   []interface{}
	}{
		{[]int{1, 2, 3}, func(i interface{}) interface{} {
			return i.(int) * 2
		}, []interface{}{2, 4, 6}},
		{"str", func(i interface{}) interface{} {
			return string(i.(rune)) + "1"
		}, []interface{}{"s1", "t1", "r1"}},
	}

	for _, test := range tests {
		if q := From(test.input).Select(test.selector); !validateQuery(q, test.output) {
			t.Errorf("From(%v).Select()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestSelectT(t *testing.T) {
	tests := []struct {
		input    interface{}
		selector interface{}
		output   []interface{}
	}{
		{[]int{1, 2, 3}, func(i interface{}) interface{} {
			return i.(int) * 2
		}, []interface{}{2, 4, 6}},
		{[]int{1, 2, 3}, func(i int) int {
			return i * 2
		}, []interface{}{2, 4, 6}},
		{"str", func(i interface{}) interface{} {
			return string(i.(rune)) + "1"
		}, []interface{}{"s1", "t1", "r1"}},
	}

	for _, test := range tests {
		if q := From(test.input).SelectT(test.selector); !validateQuery(q, test.output) {
			t.Errorf("From(%v).Select()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestSelectT_PanicWhenFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}
	}()

	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).SelectT(func(item, idx int) int { return item + 2 })
}

func TestSelectIndexed(t *testing.T) {
	tests := []struct {
		input    interface{}
		selector func(int, interface{}) interface{}
		output   []interface{}
	}{
		{[]int{1, 2, 3}, func(i int, x interface{}) interface{} {
			return x.(int) * i
		}, []interface{}{0, 2, 6}},
		{"str", func(i int, x interface{}) interface{} {
			return string(x.(rune)) + strconv.Itoa(i)
		}, []interface{}{"s0", "t1", "r2"}},
	}

	for _, test := range tests {
		if q := From(test.input).SelectIndexed(test.selector); !validateQuery(q, test.output) {
			t.Errorf("From(%v).SelectIndexed()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestSelectIndexedT(t *testing.T) {
	tests := []struct {
		input    interface{}
		selector interface{}
		output   []interface{}
	}{
		{[]int{1, 2, 3}, func(i int, x interface{}) interface{} {
			return x.(int) * i
		}, []interface{}{0, 2, 6}},
		{[]int{1, 2, 3}, func(i int, x int) int {
			return x * i
		}, []interface{}{0, 2, 6}},
		{"str", func(i int, x interface{}) interface{} {
			return string(x.(rune)) + strconv.Itoa(i)
		}, []interface{}{"s0", "t1", "r2"}},
	}

	for _, test := range tests {
		if q := From(test.input).SelectIndexedT(test.selector); !validateQuery(q, test.output) {
			t.Errorf("From(%v).SelectIndexed()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestSelectIndexedT_PanicWhenFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}
	}()

	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).SelectIndexedT(func(index string, item int) int { return item + 2 })
}
