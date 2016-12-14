package linq

import "testing"

func TestIntersect(t *testing.T) {
	input1 := []int{1, 2, 3}
	input2 := []int{1, 4, 7, 9, 12, 3}
	want := []interface{}{1, 3}

	if q := From(input1).Intersect(From(input2)); !validateQuery(q, want) {
		t.Errorf("From(%v).Intersect(%v)=%v expected %v", input1, input2, toSlice(q), want)
	}
}

func TestIntersectBy(t *testing.T) {
	input1 := []int{5, 7, 8}
	input2 := []int{1, 4, 7, 9, 12, 3}
	want := []interface{}{5, 8}

	if q := From(input1).IntersectBy(From(input2), func(i interface{}) interface{} {
		return i.(int) % 2
	}); !validateQuery(q, want) {
		t.Errorf("From(%v).IntersectBy(%v)=%v expected %v", input1, input2, toSlice(q), want)
	}
}

func TestIntersectByT_PanicWhenSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "IntersectByT: parameter [selectorFn] has a invalid function signature. Expected: 'func(T)T', actual: 'func(int,int)int'", func() {
		From([]int{5, 7, 8}).IntersectByT(From([]int{1, 4, 7, 9, 12, 3}), func(i, x int) int {
			return i % 2
		})
	})
}
