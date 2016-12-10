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

func TestGroupJoinT(t *testing.T) {
	tests := []struct {
		outer            interface{}
		inner            interface{}
		outerKeySelector interface{}
		innerKeySelector interface{}
		resultSelector   interface{}
		want             []interface{}
	}{
		{
			[]int{0, 1, 2}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			func(i int) int { return i },
			func(i int) int { return i % 2 },
			func(outer int, inners []int) KeyValue { return KeyValue{outer, len(inners)} },
			[]interface{}{KeyValue{0, 4}, KeyValue{1, 5}, KeyValue{2, 0}},
		},
	}

	for _, test := range tests {
		if q := From(test.outer).GroupJoinT(From(test.inner), test.outerKeySelector, test.innerKeySelector, test.resultSelector); !validateQuery(q, test.want) {
			t.Errorf("From().GroupJoin()=%v expected %v", toSlice(q), test.want)
		}
	}
}

func TestGroupJoinT_PanicWhenOuterKeySelectorFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	From([]int{0, 1, 2}).GroupJoinT(
		From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
		func(i, j int) int { return i },
		func(i int) int { return i % 2 },
		func(outer int, inners []int) KeyValue { return KeyValue{outer, len(inners)} },
	)
}

func TestGroupJoinT_PanicWhenInnerKeySelectorFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	From([]int{0, 1, 2}).GroupJoinT(
		From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
		func(i int) int { return i },
		func(i, j int) int { return i % 2 },
		func(outer int, inners []int) KeyValue { return KeyValue{outer, len(inners)} },
	)
}

func TestGroupJoinT_PanicWhenResultSelectorFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	From([]int{0, 1, 2}).GroupJoinT(
		From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
		func(i int) int { return i },
		func(i int) int { return i % 2 },
		func(outer, j int, inners []int) KeyValue { return KeyValue{outer, len(inners)} },
	)
}
