package linq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExcept(t *testing.T) {
	input1 := []int{1, 2, 3, 4, 5, 1, 2, 5}
	input2 := []int{1, 2}
	want := []interface{}{3, 4, 5, 5}

	if q := From(input1).Except(From(input2)); !validateQuery(q, want) {
		t.Errorf("From(%v).Except(%v)=%v expected %v", input1, input2, toSlice(q), want)
	}
}

func TestExceptG(t *testing.T) {
	input1 := []int{1, 2, 3, 4, 5, 1, 2, 5}
	input2 := []int{1, 2}
	want := []int{3, 4, 5, 5}

	actual := FromSliceG(input1).Except(FromSliceG(input2)).ToSlice()
	assert.Equal(t, want, actual)
}

func TestExceptBy(t *testing.T) {
	input1 := []int{1, 2, 3, 4, 5, 1, 2, 5}
	input2 := []int{1}
	want := []interface{}{2, 4, 2}

	if q := From(input1).ExceptBy(From(input2), func(i interface{}) interface{} {
		return i.(int) % 2
	}); !validateQuery(q, want) {
		t.Errorf("From(%v).ExceptBy(%v)=%v expected %v", input1, input2, toSlice(q), want)
	}
}

func TestExceptByG(t *testing.T) {
	input1 := []int{1, 2, 3, 4, 5, 1, 2, 5}
	input2 := []int{1}
	want := []int{2, 4, 2}

	assert.Equal(t, want, FromSliceG(input1).Expend(To2[int, int]()).(*Expended[int, int]).ExceptBy(FromSliceG(input2), func(i int) int {
		return i % 2
	}).ToSlice())
}

func TestExceptByT_PanicWhenSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "ExceptByT: parameter [selectorFn] has a invalid function signature. Expected: 'func(T)T', actual: 'func(int,int)int'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).ExceptByT(From([]int{1}), func(x, item int) int { return item + 2 })
	})
}
