package linq

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestZip(t *testing.T) {
	input1 := []int{1, 2, 3}
	input2 := []int{2, 4, 5, 1}
	want := []interface{}{3, 6, 8}

	if q := From(input1).Zip(From(input2), func(i, j interface{}) interface{} {
		return i.(int) + j.(int)
	}); !validateQuery(q, want) {
		t.Errorf("From(%v).Zip(%v)=%v expected %v", input1, input2, toSlice(q), want)
	}
}

func TestZipG(t *testing.T) {
	input1 := []int{1, 2, 3}
	input2 := []int{2, 4, 5, 1}
	want := []string{"3", "6", "8"}

	slice := FromSliceG(input1).Expend(Expend3[int, int, string]()).(Expended3[int, int, string]).Zip(FromSliceG(input2), func(i1, i2 int) string {
		return strconv.Itoa(i1 + i2)
	}).ToSlice()
	assert.Equal(t, want, slice)
}

func TestZipT_PanicWhenResultSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "ZipT: parameter [resultSelectorFn] has a invalid function signature. Expected: 'func(T,T)T', actual: 'func(int,int,int)int'", func() {
		input1 := []int{1, 2, 3}
		input2 := []int{2, 4, 5, 1}

		From(input1).ZipT(From(input2), func(i, j, k int) int {
			return i + j
		})
	})
}
