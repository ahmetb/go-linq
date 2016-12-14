package linq

import "testing"

func TestJoin(t *testing.T) {
	outer := []int{0, 1, 2, 3, 4, 5, 8}
	inner := []int{1, 2, 1, 4, 7, 6, 7, 2}
	want := []interface{}{
		KeyValue{1, 1},
		KeyValue{1, 1},
		KeyValue{2, 2},
		KeyValue{2, 2},
		KeyValue{4, 4},
	}

	q := From(outer).Join(
		From(inner),
		func(i interface{}) interface{} { return i },
		func(i interface{}) interface{} { return i },
		func(outer interface{}, inner interface{}) interface{} {
			return KeyValue{outer, inner}
		})

	if !validateQuery(q, want) {
		t.Errorf("From().Join()=%v expected %v", toSlice(q), want)
	}
}

func TestJoinT_PanicWhenOuterKeySelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "JoinT: parameter [outerKeySelectorFn] has a invalid function signature. Expected: 'func(T)T', actual: 'func(int,int)int'", func() {
		From([]int{0, 1, 2}).JoinT(
			From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			func(i, j int) int { return i },
			func(i int) int { return i % 2 },
			func(outer int, inner int) KeyValue { return KeyValue{outer, inner} },
		)
	})
}

func TestJoinT_PanicWhenInnerKeySelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "JoinT: parameter [innerKeySelectorFn] has a invalid function signature. Expected: 'func(T)T', actual: 'func(int,int)int'", func() {
		From([]int{0, 1, 2}).JoinT(
			From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			func(i int) int { return i },
			func(i, j int) int { return i % 2 },
			func(outer int, inners []int) KeyValue { return KeyValue{outer, len(inners)} },
		)
	})
}

func TestJoinT_PanicWhenResultSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "JoinT: parameter [resultSelectorFn] has a invalid function signature. Expected: 'func(T,T)T', actual: 'func(int,int,int)linq.KeyValue'", func() {
		From([]int{0, 1, 2}).JoinT(
			From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			func(i int) int { return i },
			func(i int) int { return i % 2 },
			func(outer int, inner, j int) KeyValue { return KeyValue{outer, inner} },
		)
	})
}
