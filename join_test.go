package linq

import "testing"

func TestJoin(t *testing.T) {
	outer := []int{0, 1, 2, 3, 4, 5, 8}
	inner := []int{1, 2, 1, 4, 7, 6, 7, 2}
	want := []any{
		legacyKeyValue{1, 1},
		legacyKeyValue{1, 1},
		legacyKeyValue{2, 2},
		legacyKeyValue{2, 2},
		legacyKeyValue{4, 4},
	}

	q := legacyFrom(outer).Join(
		legacyFrom(inner),
		func(i any) any { return i },
		func(i any) any { return i },
		func(outer any, inner any) any {
			return legacyKeyValue{outer, inner}
		})

	if !testQueryIteration(q, want) {
		t.Errorf("From().Join()=%v expected %v", toSlice(q), want)
	}
}

func TestJoinT_PanicWhenOuterKeySelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "JoinT: parameter [outerKeySelectorFn] has a invalid function signature. Expected: 'func(T)T', actual: 'func(int,int)int'", func() {
		legacyFrom([]int{0, 1, 2}).JoinT(
			legacyFrom([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			func(i, j int) int { return i },
			func(i int) int { return i % 2 },
			func(outer int, inner int) legacyKeyValue { return legacyKeyValue{outer, inner} },
		)
	})
}

func TestJoinT_PanicWhenInnerKeySelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "JoinT: parameter [innerKeySelectorFn] has a invalid function signature. Expected: 'func(T)T', actual: 'func(int,int)int'", func() {
		legacyFrom([]int{0, 1, 2}).JoinT(
			legacyFrom([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			func(i int) int { return i },
			func(i, j int) int { return i % 2 },
			func(outer int, inners []int) legacyKeyValue { return legacyKeyValue{outer, len(inners)} },
		)
	})
}

func TestJoinT_PanicWhenResultSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "JoinT: parameter [resultSelectorFn] has a invalid function signature. Expected: 'func(T,T)T', actual: 'func(int,int,int)linq.legacyKeyValue'", func() {
		legacyFrom([]int{0, 1, 2}).JoinT(
			legacyFrom([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			func(i int) int { return i },
			func(i int) int { return i % 2 },
			func(outer int, inner, j int) legacyKeyValue { return legacyKeyValue{outer, inner} },
		)
	})
}
