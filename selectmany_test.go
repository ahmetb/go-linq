package linq

import (
	"strconv"
	"testing"
)

func TestSelectMany(t *testing.T) {
	tests := []struct {
		input    interface{}
		selector func(interface{}) Query
		output   []interface{}
	}{
		{[][]int{{1, 2, 3}, {4, 5, 6, 7}}, func(i interface{}) Query {
			return From(i)
		}, []interface{}{1, 2, 3, 4, 5, 6, 7}},
		{[]string{"str", "ing"}, func(i interface{}) Query {
			return FromString(i.(string))
		}, []interface{}{'s', 't', 'r', 'i', 'n', 'g'}},
	}

	for _, test := range tests {
		if q := From(test.input).SelectMany(test.selector); !validateQuery(q, test.output) {
			t.Errorf("From(%v).SelectMany()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestSelectManyIndexed(t *testing.T) {
	tests := []struct {
		input    interface{}
		selector func(int, interface{}) Query
		output   []interface{}
	}{
		{[][]int{{1, 2, 3}, {4, 5, 6, 7}}, func(i int, x interface{}) Query {
			if i > 0 {
				return From(x.([]int)[1:])
			}
			return From(x)
		}, []interface{}{1, 2, 3, 5, 6, 7}},
		{[]string{"str", "ing"}, func(i int, x interface{}) Query {
			return FromString(x.(string) + strconv.Itoa(i))
		}, []interface{}{'s', 't', 'r', '0', 'i', 'n', 'g', '1'}},
	}

	for _, test := range tests {
		if q := From(test.input).SelectManyIndexed(test.selector); !validateQuery(q, test.output) {
			t.Errorf("From(%v).SelectManyIndexed()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestSelectManyBy(t *testing.T) {
	tests := []struct {
		input          interface{}
		selector       func(interface{}) Query
		resultSelector func(interface{}, interface{}) interface{}
		output         []interface{}
	}{
		{[][]int{{1, 2, 3}, {4, 5, 6, 7}}, func(i interface{}) Query {
			return From(i)
		}, func(x interface{}, y interface{}) interface{} {
			return y.(int) + 1
		}, []interface{}{2, 3, 4, 5, 6, 7, 8}},
		{[]string{"str", "ing"}, func(i interface{}) Query {
			return FromString(i.(string))
		}, func(x interface{}, y interface{}) interface{} {
			return string(y.(rune)) + "_"
		}, []interface{}{"s_", "t_", "r_", "i_", "n_", "g_"}},
	}

	for _, test := range tests {
		if q := From(test.input).SelectManyBy(test.selector, test.resultSelector); !validateQuery(q, test.output) {
			t.Errorf("From(%v).SelectManyBy()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}

func TestSelectManyIndexedBy(t *testing.T) {
	tests := []struct {
		input          interface{}
		selector       func(int, interface{}) Query
		resultSelector func(interface{}, interface{}) interface{}
		output         []interface{}
	}{
		{[][]int{{1, 2, 3}, {4, 5, 6, 7}}, func(i int, x interface{}) Query {
			if i == 0 {
				return From([]int{10, 20, 30})
			}
			return From(x)
		}, func(x interface{}, y interface{}) interface{} {
			return y.(int) + 1
		}, []interface{}{11, 21, 31, 5, 6, 7, 8}},
		{[]string{"st", "ng"}, func(i int, x interface{}) Query {
			if i == 0 {
				return FromString(x.(string) + "r")
			}
			return FromString("i" + x.(string))
		}, func(x interface{}, y interface{}) interface{} {
			return string(y.(rune)) + "_"
		}, []interface{}{"s_", "t_", "r_", "i_", "n_", "g_"}},
	}

	for _, test := range tests {
		if q := From(test.input).SelectManyByIndexed(test.selector, test.resultSelector); !validateQuery(q, test.output) {
			t.Errorf("From(%v).SelectManyIndexedBy()=%v expected %v", test.input, toSlice(q), test.output)
		}
	}
}
