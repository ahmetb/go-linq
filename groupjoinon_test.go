package linq

import (
	"testing"
)

func TestGroupJoinOn(t *testing.T) {
	outer := []int{0, 1, 2}
	inner := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	want := []interface{}{
		KeyValue{0, 9},
		KeyValue{1, 8},
		KeyValue{2, 7},
	}

	q := From(outer).GroupJoinOn(
		From(inner),
		func(i interface{}, j interface{}) bool { return i.(int) < j.(int) },
		func(outer interface{}, inners []interface{}) interface{} {
			return KeyValue{outer, len(inners)}
		})

	if !validateQuery(q, want) {
		t.Errorf("From().GroupJoinOn()=%v expected %v", toSlice(q), want)
	}
}

func TestGroupJoinOnT(t *testing.T) {
	outer := []int{0, 1, 2}
	inner := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	want := []interface{}{
		KeyValue{0, 9},
		KeyValue{1, 8},
		KeyValue{2, 7},
	}

	q := From(outer).GroupJoinOnT(
		From(inner),
		func(i, j int) bool { return i < j },
		func(outer int, inners []int) KeyValue {
			return KeyValue{outer, len(inners)}
		})

	if !validateQuery(q, want) {
		t.Errorf("From().GroupJoinOnT()=%v expected %v", toSlice(q), want)
	}
}

func TestGroupJoinOnT_PanicWhenOnPredicateFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "GroupJoinOnT: parameter [onPredicateFn] has a invalid function signature. Expected: 'func(T,T)bool', actual: 'func(int,int)int'", func() {
		From([]int{0, 1, 2}).GroupJoinOnT(
			From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			func(i, j int) int { return i + j },
			func(outer int, inners []int) KeyValue { return KeyValue{outer, len(inners)} },
		)
	})
}

func TestGroupJoinOnT_PanicWhenResultSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "GroupJoinOnT: parameter [resultSelectorFn] has a invalid function signature. Expected: 'func(T,T)T', actual: 'func(int,int,[]int)linq.KeyValue'", func() {
		From([]int{0, 1, 2}).GroupJoinOnT(
			From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			func(i, j int) bool { return i+j == 5 },
			func(outer, j int, inners []int) KeyValue { return KeyValue{outer, len(inners)} },
		)
	})
}
