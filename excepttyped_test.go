package linq

import "testing"

func TestExceptByT(t *testing.T) {
	tests := []struct {
		input1   interface{}
		input2   interface{}
		selector interface{}
		want     []interface{}
	}{
		{[]int{1, 2, 3, 4, 5, 1, 2, 5}, []int{1}, func(i interface{}) interface{} {
			return i.(int) % 2
		}, []interface{}{2, 4, 2}},
		{[]int{1, 2, 3, 4, 5, 1, 2, 5}, []int{1}, func(i int) int {
			return i % 2
		}, []interface{}{2, 4, 2}},
	}

	for _, test := range tests {
		if q := From(test.input1).ExceptByT(From(test.input2), test.selector); !validateQuery(q, test.want) {
			t.Errorf("From(%v).ExceptBy(%v)=%v expected %v", test.input1, test.input2, toSlice(q), test.want)
		}
	}
}

func TestExceptByT_PanicWhenFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).ExceptByT(From([]int{1}), func(x, item int) int { return item + 2 })
}
