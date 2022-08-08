package linq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntersect(t *testing.T) {
	input1 := []int{1, 2, 3}
	input2 := []int{1, 4, 7, 9, 12, 3}
	want := []interface{}{1, 3}

	if q := From(input1).Intersect(From(input2)); !validateQuery(q, want) {
		t.Errorf("From(%v).Intersect(%v)=%v expected %v", input1, input2, toSlice(q), want)
	}
}

func TestIntersectG(t *testing.T) {
	input1 := []int{1, 2, 3}
	input2 := []int{1, 4, 7, 9, 12, 3}
	want := []int{1, 3}

	assert.Equal(t, want, FromSliceG(input1).Intersect(FromSliceG(input2)).ToSlice())
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

func TestIntersectByG(t *testing.T) {
	input1 := []int{5, 7, 8}
	input2 := []int{1, 5, 7, 9, 12, 3}
	want := []int{5, 8}
	actual := FromSliceG(input1).IntersectBy(FromSliceG(input2), IntersectSelector(func(i int) int {
		return i % 2
	})).ToSlice()
	assert.Equal(t, want, actual)
}

func TestIntersectByT_PanicWhenSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "IntersectByT: parameter [selectorFn] has a invalid function signature. Expected: 'func(T)T', actual: 'func(int,int)int'", func() {
		From([]int{5, 7, 8}).IntersectByT(From([]int{1, 4, 7, 9, 12, 3}), func(i, x int) int {
			return i % 2
		})
	})
}
