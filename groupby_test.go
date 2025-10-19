package linq

import (
	"reflect"
	"testing"
)

func TestGroupBy(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	wantEven := []any{2, 4, 6, 8}
	wantOdd := []any{1, 3, 5, 7, 9}

	q := From(input).GroupBy(
		func(i any) any { return i.(int) % 2 },
		func(i any) any { return i.(int) },
	)

	eq := true
	for item := range q.Iterate {
		group := item.(Group)
		switch group.Key.(int) {
		case 0:
			if !reflect.DeepEqual(group.Group, wantEven) {
				eq = false
			}
		case 1:
			if !reflect.DeepEqual(group.Group, wantOdd) {
				eq = false
			}
		default:
			eq = false
		}
	}

	if !eq {
		t.Errorf("From(%v).GroupBy()=%v", input, toSlice(q))
	}
}

func TestGroupBy_Abort(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	q := From(input).GroupBy(
		func(i any) any { return i.(int) % 2 },
		func(i any) any { return i.(int) },
	)

	runDryIteration(q)
}

func TestGroupByT_PanicWhenKeySelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "GroupByT: parameter [keySelectorFn] has a invalid function signature. Expected: 'func(T)T', actual: 'func(int,int)bool'", func() {
		var r []int
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).GroupByT(
			func(i, j int) bool { return true },
			func(i int) int { return i },
		).ToSlice(&r)
	})
}

func TestGroupByT_PanicWhenElementSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "GroupByT: parameter [elementSelectorFn] has a invalid function signature. Expected: 'func(T)T', actual: 'func(int,int)int'", func() {
		var r []int
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).GroupByT(
			func(i int) bool { return true },
			func(i, j int) int { return i },
		).ToSlice(&r)
	})
}
