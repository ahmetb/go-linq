package linq

import "testing"

func TestGroupJoin(t *testing.T) {
	outer := []int{0, 1, 2}
	inner := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	want := []any{
		legacyKeyValue{0, 4},
		legacyKeyValue{1, 5},
		legacyKeyValue{2, 0},
	}

	q := legacyFrom(outer).GroupJoin(
		legacyFrom(inner),
		func(i any) any { return i },
		func(i any) any { return i.(int) % 2 },
		func(outer any, inners []any) any {
			return legacyKeyValue{outer, len(inners)}
		})

	if !testQueryIteration(q, want) {
		t.Errorf("From().GroupJoin()=%v expected %v", toSlice(q), want)
	}
}

func TestGroupJoinT_PanicWhenOuterKeySelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "GroupJoinT: parameter [outerKeySelectorFn] has a invalid function signature. Expected: 'func(T)T', actual: 'func(int,int)int'", func() {
		legacyFrom([]int{0, 1, 2}).GroupJoinT(
			legacyFrom([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			func(i, j int) int { return i },
			func(i int) int { return i % 2 },
			func(outer int, inners []int) legacyKeyValue { return legacyKeyValue{outer, len(inners)} },
		)
	})
}

func TestGroupJoinT_PanicWhenInnerKeySelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "GroupJoinT: parameter [innerKeySelectorFn] has a invalid function signature. Expected: 'func(T)T', actual: 'func(int,int)int'", func() {
		legacyFrom([]int{0, 1, 2}).GroupJoinT(
			legacyFrom([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			func(i int) int { return i },
			func(i, j int) int { return i % 2 },
			func(outer int, inners []int) legacyKeyValue { return legacyKeyValue{outer, len(inners)} },
		)
	})
}

func TestGroupJoinT_PanicWhenResultSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "GroupJoinT: parameter [resultSelectorFn] has a invalid function signature. Expected: 'func(T,T)T', actual: 'func(int,int,[]int)linq.legacyKeyValue'", func() {
		legacyFrom([]int{0, 1, 2}).GroupJoinT(
			legacyFrom([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			func(i int) int { return i },
			func(i int) int { return i % 2 },
			func(outer, j int, inners []int) legacyKeyValue { return legacyKeyValue{outer, len(inners)} },
		)
	})
}
