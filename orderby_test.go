package linq

import (
	"iter"
	"testing"
)

func TestEmpty(t *testing.T) {
	q := From([]string{}).OrderBy(func(in any) any {
		return 0
	})

	next, stop := iter.Pull(q.Iterate)
	defer stop()

	_, ok := next()
	if ok {
		t.Errorf("Iterator for empty collection must return ok=false")
	}
}

func TestOrderBy(t *testing.T) {
	slice := make([]foo, 100)

	for i := len(slice) - 1; i >= 0; i-- {
		slice[i].f1 = i
	}

	q := From(slice).OrderBy(func(i any) any {
		return i.(foo).f1
	})

	next, stop := iter.Pull(q.Iterate)
	defer stop()

	j := 0
	for item, ok := next(); ok; item, ok = next() {
		if item.(foo).f1 != j {
			t.Errorf("OrderBy()[%v]=%v expected %v", j, item, foo{f1: j})
		}

		j++
	}
}

func TestOrderByT_PanicWhenSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "OrderByT: parameter [selectorFn] has a invalid function signature. Expected: 'func(T)T', actual: 'func(int,int)int'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).OrderByT(func(item, j int) int { return item + 2 })
	})
}

func TestOrderByDescending(t *testing.T) {
	slice := make([]foo, 100)

	for i := 0; i < len(slice); i++ {
		slice[i].f1 = i
	}

	q := From(slice).OrderByDescending(func(i any) any {
		return i.(foo).f1
	})

	next, stop := iter.Pull(q.Iterate)
	defer stop()

	j := len(slice) - 1
	for item, ok := next(); ok; item, ok = next() {
		if item.(foo).f1 != j {
			t.Errorf("OrderByDescending()[%v]=%v expected %v", j, item, foo{f1: j})
		}

		j--
	}
}

func TestOrderByDescendingT_PanicWhenSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "OrderByDescendingT: parameter [selectorFn] has a invalid function signature. Expected: 'func(T)T', actual: 'func(int,int)int'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).OrderByDescendingT(func(item, j int) int { return item + 2 })
	})
}

func TestThenBy(t *testing.T) {
	slice := make([]foo, 1000)

	for i := len(slice) - 1; i >= 0; i-- {
		slice[i].f1 = i
		slice[i].f2 = i%2 == 0
	}

	q := From(slice).OrderBy(func(i any) any {
		return i.(foo).f2
	}).ThenBy(func(i any) any {
		return i.(foo).f1
	})

	next, stop := iter.Pull(q.Iterate)
	defer stop()

	for item, ok := next(); ok; item, ok = next() {
		if item.(foo).f2 != (item.(foo).f1%2 == 0) {
			t.Errorf("OrderBy().ThenBy()=%v", item)
		}
	}
}

func TestThenBy_Abort(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	q := From(input).OrderBy(func(i any) any {
		return i.(int)
	}).ThenBy(func(i any) any {
		return i.(int)
	})

	runDryIteration(q.Query)
}

func TestThenByT_PanicWhenSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "ThenByT: parameter [selectorFn] has a invalid function signature. Expected: 'func(T)T', actual: 'func(int,int)bool'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).
			OrderByT(func(item int) int { return item }).
			ThenByT(func(item, j int) bool { return true })
	})
}

func TestThenByDescending(t *testing.T) {
	slice := make([]foo, 1000)

	for i := len(slice) - 1; i >= 0; i-- {
		slice[i].f1 = i
		slice[i].f2 = i%2 == 0
	}

	q := From(slice).OrderBy(func(i any) any {
		return i.(foo).f2
	}).ThenByDescending(func(i any) any {
		return i.(foo).f1
	})

	next, stop := iter.Pull(q.Iterate)
	defer stop()

	for item, ok := next(); ok; item, ok = next() {
		if item.(foo).f2 != (item.(foo).f1%2 == 0) {
			t.Errorf("OrderBy().ThenByDescending()=%v", item)
		}
	}
}

func TestThenByDescending_Abort(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	q := From(input).OrderBy(func(i any) any {
		return i.(int)
	}).ThenByDescending(func(i any) any {
		return i.(int)
	})

	runDryIteration(q.Query)
}

func TestThenByDescendingT_PanicWhenSelectorFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "ThenByDescending: parameter [selectorFn] has a invalid function signature. Expected: 'func(T)T', actual: 'func(int,int)bool'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).
			OrderByT(func(item int) int { return item }).
			ThenByDescendingT(func(item, j int) bool { return true })
	})
}

func TestSort(t *testing.T) {
	slice := make([]foo, 100)

	for i := len(slice) - 1; i >= 0; i-- {
		slice[i].f1 = i
	}

	q := From(slice).Sort(func(i, j any) bool {
		return i.(foo).f1 < j.(foo).f1
	})

	next, stop := iter.Pull(q.Iterate)
	defer stop()

	j := 0
	for item, ok := next(); ok; item, ok = next() {
		if item.(foo).f1 != j {
			t.Errorf("Sort()[%v]=%v expected %v", j, item, foo{f1: j})
		}

		j++
	}
}

func TestSort_Abort(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	q := From(input).Sort(func(i, j any) bool {
		return i.(int) < j.(int)
	})

	runDryIteration(q)
}

func TestSortT_PanicWhenLessFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "SortT: parameter [lessFn] has a invalid function signature. Expected: 'func(T,T)bool', actual: 'func(int,int)string'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).SortT(func(i, j int) string { return "" })
	})
}
