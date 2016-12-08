package linq

import "testing"

func TestJoinT(t *testing.T) {
	outer := []int{0, 1, 2, 3, 4, 5, 8}
	inner := []int{1, 2, 1, 4, 7, 6, 7, 2}
	want := []interface{}{
		KeyValue{1, 1},
		KeyValue{1, 1},
		KeyValue{2, 2},
		KeyValue{2, 2},
		KeyValue{4, 4},
	}

	q := From(outer).JoinT(
		From(inner),
		func(i interface{}) interface{} { return i },
		func(i interface{}) interface{} { return i },
		func(outer interface{}, inner interface{}) interface{} {
			return KeyValue{outer, inner}
		})

	if !validateQuery(q, want) {
		t.Errorf("From().Join()=%v expected %v", toSlice(q), want)
	}

	tests := []struct {
		outer            interface{}
		inner            interface{}
		outerKeySelector interface{}
		innerKeySelector interface{}
		resultSelector   interface{}
		want             []interface{}
	}{
		{
			[]int{0, 1, 2, 3, 4, 5, 8}, []int{1, 2, 1, 4, 7, 6, 7, 2},
			func(i interface{}) interface{} { return i },
			func(i interface{}) interface{} { return i },
			func(outer interface{}, inner interface{}) interface{} { return KeyValue{outer, inner} },
			[]interface{}{KeyValue{1, 1}, KeyValue{1, 1}, KeyValue{2, 2}, KeyValue{2, 2}, KeyValue{4, 4}},
		},
		{
			[]int{0, 1, 2, 3, 4, 5, 8}, []int{1, 2, 1, 4, 7, 6, 7, 2},
			func(i int) int { return i },
			func(i int) int { return i },
			func(outer int, inner int) KeyValue { return KeyValue{outer, inner} },
			[]interface{}{KeyValue{1, 1}, KeyValue{1, 1}, KeyValue{2, 2}, KeyValue{2, 2}, KeyValue{4, 4}},
		},
	}

	for _, test := range tests {
		if q := From(test.outer).JoinT(From(test.inner), test.outerKeySelector, test.innerKeySelector, test.resultSelector); !validateQuery(q, test.want) {
			t.Errorf("From().Join()=%v expected %v", toSlice(q), test.want)
		}
	}
}

func TestJoinT_PanicWhenOuterKeySelectorFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	From([]int{0, 1, 2}).JoinT(
		From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
		func(i, j int) int { return i },
		func(i int) int { return i % 2 },
		func(outer int, inner int) KeyValue { return KeyValue{outer, inner} },
	)
}

func TestJoinT_PanicWhenInnerKeySelectorFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	From([]int{0, 1, 2}).JoinT(
		From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
		func(i int) int { return i },
		func(i, j int) int { return i % 2 },
		func(outer int, inners []int) KeyValue { return KeyValue{outer, len(inners)} },
	)
}

func TestJoinT_PanicWhenResultSelectorFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	From([]int{0, 1, 2}).JoinT(
		From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
		func(i int) int { return i },
		func(i int) int { return i % 2 },
		func(outer int, inner, j int) KeyValue { return KeyValue{outer, inner} },
	)
}
