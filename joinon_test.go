package linq

import (
	"testing"
)

func TestJoinOn(t *testing.T) {
	outer := []int{0, 1}
	inner := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	want := []interface{}{
		KeyValue{0, 5},
		KeyValue{1, 4},
	}

	q := From(outer).JoinOn(
		From(inner),
		func(i interface{}, j interface{}) bool { return i.(int)+j.(int) == 5 },
		func(outer interface{}, inner interface{}) interface{} {
			return KeyValue{outer, inner}
		})

	if !validateQuery(q, want) {
		t.Errorf("From().JoinOn()=%v expected %v", toSlice(q), want)
	}
}

func TestJoinOnT(t *testing.T) {
	outer := []int{0, 1}
	inner := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	want := []interface{}{
		KeyValue{0, 5},
		KeyValue{1, 4},
	}

	q := From(outer).JoinOnT(
		From(inner),
		func(i, j int) bool { return i+j == 5 },
		func(outer int, inner int) KeyValue {
			return KeyValue{outer, inner}
		})

	if !validateQuery(q, want) {
		t.Errorf("From().JoinOnT()=%v expected %v", toSlice(q), want)
	}
}

func TestJoinOnT_PanicWhenOnPredicateFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "JoinOnT: parameter [onPredicateFn] has a invalid function signature. Expected: 'func(T,T)bool', actual: 'func(int,int)int'", func() {
		From([]int{0, 1}).JoinOnT(
			From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			func(i, j int) int { return i + j },
			func(outer int, inner int) KeyValue { return KeyValue{outer, inner} },
		)
	})
}

func TestJoinOnT_PanicWhenResultSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "JoinOnT: parameter [resultSelectorFn] has a invalid function signature. Expected: 'func(T,T)T', actual: 'func(int,int,int)linq.KeyValue'", func() {
		From([]int{0, 1, 2}).JoinOnT(
			From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			func(i, j int) bool { return i == j%2 },
			func(outer int, inner, j int) KeyValue { return KeyValue{outer, inner} },
		)
	})
}
