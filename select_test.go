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
