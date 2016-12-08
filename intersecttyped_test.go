package linq

import "testing"

func TestIntersectByT(t *testing.T) {
	tests := []struct {
		input1   interface{}
		input2   interface{}
		selector interface{}
		want     []interface{}
	}{
		{
			[]int{5, 7, 8}, []int{1, 4, 7, 9, 12, 3}, func(i interface{}) interface{} {
				return i.(int) % 2
			}, []interface{}{5, 8},
		},

		{
			[]int{5, 7, 8}, []int{1, 4, 7, 9, 12, 3}, func(i int) int {
				return i % 2
			}, []interface{}{5, 8},
		},
	}

	for _, test := range tests {
		if q := From(test.input1).IntersectByT(From(test.input2), test.selector); !validateQuery(q, test.want) {
			t.Errorf("From(%v).IntersectBy(%v)=%v expected %v", test.input1, test.input2, toSlice(q), test.want)
		}
	}
}

func TestIntersectByT_PanicWhenFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	From([]int{5, 7, 8}).IntersectByT(From([]int{1, 4, 7, 9, 12, 3}), func(i, x int) int {
		return i % 2
	})
}
