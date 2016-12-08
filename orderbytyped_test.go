package linq

import "testing"

func TestOrderByT(t *testing.T) {
	slice := make([]foo, 100)

	for i := len(slice) - 1; i >= 0; i-- {
		slice[i].f1 = i
	}

	q := From(slice).OrderByT(func(i interface{}) interface{} {
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

func TestOrderByT_PanicWhenFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).OrderByT(func(item, j int) int { return item + 2 })
}

func TestOrderByDescendingT(t *testing.T) {
	slice := make([]foo, 100)

	for i := 0; i < len(slice); i++ {
		slice[i].f1 = i
	}

	q := From(slice).OrderByDescendingT(func(i foo) int {
		return i.f1
	})

	j := len(slice) - 1
	next := q.Iterate()
	for item, ok := next(); ok; item, ok = next() {
		if item.(foo).f1 != j {
			t.Errorf("OrderByDescendingT()[%v]=%v expected %v", j, item, foo{f1: j})
		}

		j--
	}
}

func TestOrderByDescendingT_PanicWhenFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).OrderByDescendingT(func(item, j int) int { return item + 2 })
}

func TestThenByT(t *testing.T) {
	slice := make([]foo, 1000)

	for i := len(slice) - 1; i >= 0; i-- {
		slice[i].f1 = i
		slice[i].f2 = i%2 == 0
	}

	q := From(slice).OrderByT(func(i foo) bool {
		return i.f2
	}).ThenByT(func(i foo) int {
		return i.f1
	})

	next := q.Iterate()
	for item, ok := next(); ok; item, ok = next() {
		if item.(foo).f2 != (item.(foo).f1%2 == 0) {
			t.Errorf("OrderBy().ThenBy()=%v", item)
		}
	}
}

func TestThenByT_PanicWhenFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).
		OrderByT(func(item int) int { return item }).
		ThenByT(func(item, j int) bool { return true })
}

func TestThenByTDescending(t *testing.T) {
	slice := make([]foo, 1000)

	for i := len(slice) - 1; i >= 0; i-- {
		slice[i].f1 = i
		slice[i].f2 = i%2 == 0
	}

	q := From(slice).OrderByT(func(i foo) bool {
		return i.f2
	}).ThenByDescendingT(func(i foo) int {
		return i.f1
	})

	next := q.Iterate()
	for item, ok := next(); ok; item, ok = next() {
		if item.(foo).f2 != (item.(foo).f1%2 == 0) {
			t.Errorf("OrderByT().ThenByDescendingT()=%v", item)
		}
	}
}

func TestThenByDescendingT_PanicWhenFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).
		OrderByT(func(item int) int { return item }).
		ThenByDescendingT(func(item, j int) bool { return true })
}

func TestSortT(t *testing.T) {
	slice := make([]foo, 100)

	for i := len(slice) - 1; i >= 0; i-- {
		slice[i].f1 = i
	}

	q := From(slice).SortT(func(i, j foo) bool {
		return i.f1 < j.f1
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

func TestSortT_PanicWhenFunctionIsNotValid(t *testing.T) {
	defer func() {
		r := recover()
		t.Log(r)
		if r == nil {
			t.Error("This execution should panic", r)
		}

	}()

	From([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}).SortT(func(i, j int) string { return "" })
}
