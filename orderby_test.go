package linq

import (
	"iter"
	"testing"
)

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func TestEmpty(t *testing.T) {
	q := FromSlice([]string{}).OrderBy(func(in string) int {
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

	q := FromSlice(slice).OrderBy(func(f foo) int {
		return f.f1
	})

	next, stop := iter.Pull(q.Iterate)
	defer stop()

	j := 0
	for item, ok := next(); ok; item, ok = next() {
		if item.f1 != j {
			t.Errorf("OrderBy()[%v]=%v expected %v", j, item, foo{f1: j})
		}

		j++
	}
}

func TestOrderByDescending(t *testing.T) {
	slice := make([]foo, 100)

	for i := 0; i < len(slice); i++ {
		slice[i].f1 = i
	}

	q := FromSlice(slice).OrderByDescending(func(f foo) int {
		return f.f1
	})

	next, stop := iter.Pull(q.Iterate)
	defer stop()

	j := len(slice) - 1
	for item, ok := next(); ok; item, ok = next() {
		if item.f1 != j {
			t.Errorf("OrderByDescending()[%v]=%v expected %v", j, item, foo{f1: j})
		}

		j--
	}
}

func TestThenBy(t *testing.T) {
	slice := make([]foo, 1000)

	for i := len(slice) - 1; i >= 0; i-- {
		slice[i].f1 = i
		slice[i].f2 = i%2 == 0
	}

	q := FromSlice(slice).OrderBy(func(f foo) int {
		return boolToInt(f.f2)
	}).ThenBy(func(f foo) int {
		return f.f1
	})

	next, stop := iter.Pull(q.Iterate)
	defer stop()

	prevByGroup := map[bool]int{true: -1, false: -1}
	for item, ok := next(); ok; item, ok = next() {
		if item.f2 != (item.f1%2 == 0) {
			t.Errorf("OrderBy().ThenBy()=%v", item)
		}
		if item.f1 < prevByGroup[item.f2] {
			t.Errorf("OrderBy().ThenBy() not sorted by f1 within group: %v after %v", item.f1, prevByGroup[item.f2])
		}
		prevByGroup[item.f2] = item.f1
	}
}

func TestThenBy_DifferentKeyType(t *testing.T) {
	// ThenBy key type may differ from the OrderBy key type.
	slice := []foo{
		{f1: 2, f3: "b"},
		{f1: 1, f3: "b"},
		{f1: 1, f3: "a"},
	}

	q := FromSlice(slice).OrderBy(func(f foo) int {
		return f.f1
	}).ThenBy(func(f foo) string {
		return f.f3
	})

	want := []foo{
		{f1: 1, f3: "a"},
		{f1: 1, f3: "b"},
		{f1: 2, f3: "b"},
	}

	if !testQueryIteration(q.Query, want) {
		t.Errorf("OrderBy().ThenBy()=%v expected %v", toSlice(q.Query), want)
	}
}

func TestThenBy_Abort(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	q := FromSlice(input).OrderBy(func(i int) int {
		return i
	}).ThenBy(func(i int) int {
		return i
	})

	runDryIteration(q.Query)
}

func TestThenByDescending(t *testing.T) {
	slice := make([]foo, 1000)

	for i := len(slice) - 1; i >= 0; i-- {
		slice[i].f1 = i
		slice[i].f2 = i%2 == 0
	}

	q := FromSlice(slice).OrderBy(func(f foo) int {
		return boolToInt(f.f2)
	}).ThenByDescending(func(f foo) int {
		return f.f1
	})

	next, stop := iter.Pull(q.Iterate)
	defer stop()

	for item, ok := next(); ok; item, ok = next() {
		if item.f2 != (item.f1%2 == 0) {
			t.Errorf("OrderBy().ThenByDescending()=%v", item)
		}
	}
}

func TestThenByDescending_Abort(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	q := FromSlice(input).OrderBy(func(i int) int {
		return i
	}).ThenByDescending(func(i int) int {
		return i
	})

	runDryIteration(q.Query)
}

func TestSort(t *testing.T) {
	slice := make([]foo, 100)

	for i := len(slice) - 1; i >= 0; i-- {
		slice[i].f1 = i
	}

	q := FromSlice(slice).Sort(func(i, j foo) bool {
		return i.f1 < j.f1
	})

	next, stop := iter.Pull(q.Iterate)
	defer stop()

	j := 0
	for item, ok := next(); ok; item, ok = next() {
		if item.f1 != j {
			t.Errorf("Sort()[%v]=%v expected %v", j, item, foo{f1: j})
		}

		j++
	}
}

func TestSort_Abort(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	q := FromSlice(input).Sort(func(i, j int) bool {
		return i < j
	})

	runDryIteration(q)
}

// TestOrderedQueryMethodPromotion verifies that Query[T] methods, including
// generic methods like Select, are promoted through the embedded Query[T]
// field of OrderedQuery[T].
func TestOrderedQueryMethodPromotion(t *testing.T) {
	input := []int{3, 1, 2}
	want := []string{"1x", "2x", "3x"}

	q := FromSlice(input).OrderBy(func(i int) int {
		return i
	}).Select(func(i int) string {
		return string(rune('0'+i)) + "x"
	})

	if !testQueryIteration(q, want) {
		t.Errorf("OrderBy().Select()=%v expected %v", toSlice(q), want)
	}
}
