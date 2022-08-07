package linq

import (
	"github.com/stretchr/testify/assert"
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

func TestSelectGFunc(t *testing.T) {
	input := []int{1, 2, 3}
	expected := []string{"1", "2", "3"}
	stringSlice := Select[int, string](FromSliceG(input), func(i int) string {
		return strconv.Itoa(i)
	}).ToSlice()
	assert.Equal(t, expected, stringSlice)
}

func TestSelectG(t *testing.T) {
	input := []int{1, 2, 3}
	expected := []string{"1", "2", "3"}
	stringSlice := FromSliceG(input).Expend(To2[int, string]()).(Expended[int, string]).Select(func(i int) string {
		return strconv.Itoa(i)
	}).ToSlice()
	assert.Equal(t, expected, stringSlice)
}

func TestSelectIndexedG(t *testing.T) {
	input := []int{0, 1, 2}
	expected := []string{"0", "1", "2"}
	stringSlice := FromSliceG(input).Expend(To2[int, string]()).(Expended[int, string]).SelectIndexed(func(index, i int) string {
		assert.Equal(t, index, i)
		return strconv.Itoa(i)
	}).ToSlice()
	assert.Equal(t, expected, stringSlice)
}

func TestSelectT_PanicWhenSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "SelectT: parameter [selectorFn] has a invalid function signature. Expected: 'func(T)T', actual: 'func(int,int)int'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).SelectT(func(item, idx int) int { return item + 2 })
	})
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

func TestSelectIndexedT_PanicWhenSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "SelectIndexedT: parameter [selectorFn] has a invalid function signature. Expected: 'func(int,T)T', actual: 'func(string,int)int'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).SelectIndexedT(func(index string, item int) int { return item + 2 })
	})
}
