package linq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmpty(t *testing.T) {
	q := From([]string{}).OrderBy(func(in interface{}) interface{} {
		return 0
	})

	_, ok := q.Iterate()()
	if ok {
		t.Errorf("Iterator for empty collection must return ok=false")
	}
}

func TestEmptyG(t *testing.T) {
	actual := FromSliceG([]string{}).Expend(To2[string, int]()).(*Expended[string, int]).OrderBy(func(in string) int {
		return 0
	}).ToSlice()

	assert.Empty(t, actual)
}

func TestOrderBy(t *testing.T) {
	slice := make([]foo, 100)

	for i := len(slice) - 1; i >= 0; i-- {
		slice[i].f1 = i
	}

	q := From(slice).OrderBy(func(i interface{}) interface{} {
		return i.(foo).f1
	})

	j := 0
	next := q.Iterate()
	for item, ok := next(); ok; item, ok = next() {
		if item.(foo).f1 != j {
			t.Errorf("OrderBy()[%v]=%v expected %v", j, item, foo{f1: j})
		}

		j++
	}
}

func TestOrderByG(t *testing.T) {
	slice := make([]foo, 100)

	for i := len(slice) - 1; i >= 0; i-- {
		slice[i].f1 = i
	}

	actual := FromSliceG(slice).Expend(To2[foo, int]()).(*Expended[foo, int]).OrderBy(func(i foo) int {
		return i.f1
	})

	j := 0
	next := actual.Iterate()
	for item, ok := next(); ok; item, ok = next() {
		if item.f1 != j {
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

	q := From(slice).OrderByDescending(func(i interface{}) interface{} {
		return i.(foo).f1
	})

	j := len(slice) - 1
	next := q.Iterate()
	for item, ok := next(); ok; item, ok = next() {
		if item.(foo).f1 != j {
			t.Errorf("OrderByDescending()[%v]=%v expected %v", j, item, foo{f1: j})
		}

		j--
	}
}

func TestOrderByDescendingG(t *testing.T) {
	slice := make([]foo, 100)

	for i := 0; i < len(slice); i++ {
		slice[i].f1 = i
	}

	q := FromSliceG(slice).Expend(To2[foo, int]()).(*Expended[foo, int]).OrderByDescending(func(i foo) int {
		return i.f1
	})

	j := len(slice) - 1
	next := q.Iterate()
	for item, ok := next(); ok; item, ok = next() {
		if item.f1 != j {
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

	q := From(slice).OrderBy(func(i interface{}) interface{} {
		return i.(foo).f2
	}).ThenBy(func(i interface{}) interface{} {
		return i.(foo).f1
	})

	next := q.Iterate()
	for item, ok := next(); ok; item, ok = next() {
		if item.(foo).f2 != (item.(foo).f1%2 == 0) {
			t.Errorf("OrderBy().ThenBy()=%v", item)
		}
	}
}

func TestThenByG(t *testing.T) {
	slice := make([]foo, 1000)

	for i := len(slice) - 1; i >= 0; i-- {
		slice[i].f1 = i
		slice[i].f2 = i%2 == 0
	}

	q := FromSliceG(slice).Expend(To2[foo, bool]()).(*Expended[foo, bool]).OrderBy(func(i foo) bool {
		return i.f2
	}).Expend(OrderedTo2[foo, int]()).(OrderedExpended[foo, int]).ThenBy(func(i foo) int {
		return i.f1
	})

	next := q.Iterate()
	for item, ok := next(); ok; item, ok = next() {
		if item.f2 != (item.f1%2 == 0) {
			t.Errorf("OrderBy().ThenBy()=%v", item)
		}
	}
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

	q := From(slice).OrderBy(func(i interface{}) interface{} {
		return i.(foo).f2
	}).ThenByDescending(func(i interface{}) interface{} {
		return i.(foo).f1
	})

	next := q.Iterate()
	for item, ok := next(); ok; item, ok = next() {
		if item.(foo).f2 != (item.(foo).f1%2 == 0) {
			t.Errorf("OrderBy().ThenByDescending()=%v", item)
		}
	}
}

func TestThenByDescendingG(t *testing.T) {
	slice := make([]foo, 1000)

	for i := len(slice) - 1; i >= 0; i-- {
		slice[i].f1 = i
		slice[i].f2 = i%2 == 0
	}

	q := FromSliceG(slice).Expend(To2[foo, bool]()).(*Expended[foo, bool]).OrderBy(func(i foo) bool {
		return i.f2
	}).Expend(OrderedTo2[foo, int]()).(OrderedExpended[foo, int]).ThenByDescending(func(i foo) int {
		return i.f1
	})

	next := q.Iterate()
	for item, ok := next(); ok; item, ok = next() {
		if item.f2 != (item.f1%2 == 0) {
			t.Errorf("OrderBy().ThenByDescending()=%v", item)
		}
	}
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

	q := From(slice).Sort(func(i, j interface{}) bool {
		return i.(foo).f1 < j.(foo).f1
	})

	j := 0
	next := q.Iterate()
	for item, ok := next(); ok; item, ok = next() {
		if item.(foo).f1 != j {
			t.Errorf("Sort()[%v]=%v expected %v", j, item, foo{f1: j})
		}

		j++
	}
}

func TestSortG(t *testing.T) {
	slice := make([]foo, 100)

	for i := len(slice) - 1; i >= 0; i-- {
		slice[i].f1 = i
	}

	q := FromSliceG(slice).Sort(func(i, j foo) bool {
		return i.f1 < j.f1
	})

	j := 0
	next := q.Iterate()
	for item, ok := next(); ok; item, ok = next() {
		if item.f1 != j {
			t.Errorf("Sort()[%v]=%v expected %v", j, item, foo{f1: j})
		}

		j++
	}
}

func TestSortT_PanicWhenLessFnIsInvalid(t *testing.T) {
	mustPanicWithError(t, "SortT: parameter [lessFn] has a invalid function signature. Expected: 'func(T,T)bool', actual: 'func(int,int)string'", func() {
		From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).SortT(func(i, j int) string { return "" })
	})
}
