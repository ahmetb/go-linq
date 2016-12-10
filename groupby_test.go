package linq

import (
	"reflect"
	"testing"
)

func TestGroupBy(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	wantEven := []interface{}{2, 4, 6, 8}
	wantOdd := []interface{}{1, 3, 5, 7, 9}

	q := From(input).GroupBy(
		func(i interface{}) interface{} { return i.(int) % 2 },
		func(i interface{}) interface{} { return i.(int) },
	)

	next := q.Iterate()
	eq := true
	for item, ok := next(); ok; item, ok = next() {
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

func TestGroupByT(t *testing.T) {

	tests := []struct {
		input           interface{}
		keySelector     interface{}
		elementSelector interface{}
		wantEven        []interface{}
		wantOdd         []interface{}
	}{
		{
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			func(i interface{}) interface{} { return i.(int) % 2 },
			func(i interface{}) interface{} { return i.(int) },
			[]interface{}{2, 4, 6, 8}, []interface{}{1, 3, 5, 7, 9},
		},
		{
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			func(i int) int { return i % 2 },
			func(i int) int { return i },
			[]interface{}{2, 4, 6, 8}, []interface{}{1, 3, 5, 7, 9},
		},
	}
	for _, test := range tests {
		q := From(test.input).GroupByT(
			test.keySelector,
			test.elementSelector,
		)

		next := q.Iterate()
		eq := true
		for item, ok := next(); ok; item, ok = next() {
			group := item.(Group)
			switch group.Key.(int) {
			case 0:
				if !reflect.DeepEqual(group.Group, test.wantEven) {
					eq = false
				}
			case 1:
				if !reflect.DeepEqual(group.Group, test.wantOdd) {
					eq = false
				}
			default:
				eq = false
			}
		}

		if !eq {
			t.Errorf("From(%v).GroupByT()=%v", test.input, toSlice(q))
		}
	}
}

func TestGroupByT_PanicWhenKeyFunctionParameterInTypeDoesntMatch(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	var r []int
	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).GroupByT(
		func(i, j int) bool { return true },
		func(i int) int { return i },
	).ToSlice(&r)
}

func TestGroupByT_PanicWhenElementSelectorFunctionParameterInTypeDoesntMatch(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	var r []int
	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).GroupByT(
		func(i int) bool { return true },
		func(i, j int) int { return i },
	).ToSlice(&r)
}
