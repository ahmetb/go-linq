package linq

import "testing"

func TestGroupJoin(t *testing.T) {
	outer := []int{0, 1, 2}
	inner := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	want := []interface{}{
		KeyValue{0, 4},
		KeyValue{1, 5},
		KeyValue{2, 0},
	}

	q := From(outer).GroupJoin(
		From(inner),
		func(i interface{}) interface{} { return i },
		func(i interface{}) interface{} { return i.(int) % 2 },
		func(outer interface{}, inners []interface{}) interface{} {
			return KeyValue{outer, len(inners)}
		})

	if !validateQuery(q, want) {
		t.Errorf("From().GroupJoin()=%v expected %v", toSlice(q), want)
	}
}

func TestGroupJoinT_PanicWhenOuterKeySelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "GroupJoinT: parameter [outerKeySelectorFn] has a invalid function signature. Expected: 'func(T)T', actual: 'func(int,int)int'", func() {
		From([]int{0, 1, 2}).GroupJoinT(
			From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			func(i, j int) int { return i },
			func(i int) int { return i % 2 },
			func(outer int, inners []int) KeyValue { return KeyValue{outer, len(inners)} },
		)
	})
}

func TestGroupJoinT_PanicWhenInnerKeySelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "GroupJoinT: parameter [innerKeySelectorFn] has a invalid function signature. Expected: 'func(T)T', actual: 'func(int,int)int'", func() {
		From([]int{0, 1, 2}).GroupJoinT(
			From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			func(i int) int { return i },
			func(i, j int) int { return i % 2 },
			func(outer int, inners []int) KeyValue { return KeyValue{outer, len(inners)} },
		)
	})
}

func TestGroupJoinT_PanicWhenResultSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "GroupJoinT: parameter [resultSelectorFn] has a invalid function signature. Expected: 'func(T,T)T', actual: 'func(int,int,[]int)linq.KeyValue'", func() {
		From([]int{0, 1, 2}).GroupJoinT(
			From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			func(i int) int { return i },
			func(i int) int { return i % 2 },
			func(outer, j int, inners []int) KeyValue { return KeyValue{outer, len(inners)} },
		)
	})
}
