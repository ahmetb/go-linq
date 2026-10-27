package linq

import (
	"cmp"
	"iter"
	"math"
	"slices"
	"sync"
	"testing"
)

func TestOrderingCachesKeysPerIteration(t *testing.T) {
	type entry struct {
		primary   int
		secondary string
		id        int
	}
	input := []entry{{2, "b", 0}, {1, "b", 1}, {1, "a", 2}, {1, "a", 3}}
	primaryCalls, secondaryCalls := make([]int, len(input)), make([]int, len(input))
	q := FromSlice(input).OrderBy(func(e entry) int {
		primaryCalls[e.id]++
		return e.primary
	}).ThenByDescending(func(e entry) string {
		secondaryCalls[e.id]++
		return e.secondary
	})
	if primaryCalls[0] != 0 || secondaryCalls[0] != 0 {
		t.Errorf("selectors called before iteration")
	}
	for run := 1; run <= 2; run++ {
		got := q.ToSlice()
		want := []entry{input[1], input[2], input[3], input[0]}
		if !slices.Equal(got, want) {
			t.Errorf("iteration %d: %v want %v", run, got, want)
		}
		for i := range input {
			if primaryCalls[i] != run || secondaryCalls[i] != run {
				t.Errorf("item %d: primary=%d secondary=%d want %d each", i, primaryCalls[i], secondaryCalls[i], run)
			}
		}
	}
	input[0].primary = 0
	first, ok := q.First()
	if !ok || first.id != 0 {
		t.Errorf("cached keys survived into a later iteration")
	}
	for i := range input {
		if primaryCalls[i] != 3 || secondaryCalls[i] != 3 {
			t.Errorf("early consumption did not compute each key exactly once")
		}
	}
}

func TestOrderingMixedDirectionsAndNaN(t *testing.T) {
	type entry struct {
		key    float64
		second int
		id     int
	}
	input := make([]entry, 128)
	for i := range input {
		key := float64(i % 5)
		if i%7 == 0 {
			key = math.NaN()
		}
		input[i] = entry{key, i % 3, i}
	}
	for _, descending := range []bool{false, true} {
		keyCalls := 0
		selector := func(e entry) float64 { keyCalls++; return e.key }
		var q OrderedQuery[entry]
		if descending {
			q = FromSlice(input).OrderByDescending(selector)
		} else {
			q = FromSlice(input).OrderBy(selector)
		}
		q = q.ThenBy(func(e entry) int { return e.second })
		want := slices.Clone(input)
		slices.SortStableFunc(want, func(a, b entry) int {
			c := cmp.Compare(a.key, b.key)
			if descending {
				c = -c
			}
			if c != 0 {
				return c
			}
			return cmp.Compare(a.second, b.second)
		})
		got := q.ToSlice()
		if !slices.EqualFunc(got, want, func(a, b entry) bool { return a.id == b.id }) {
			t.Errorf("descending=%v: incorrect or unstable order", descending)
		}
		if keyCalls != len(input) {
			t.Errorf("selector called %d times", keyCalls)
		}
		for i, item := range input {
			if item.id != i {
				t.Errorf("ordering mutated source")
			}
		}
	}
}

func TestOrderingBranchesAndConcurrentIteration(t *testing.T) {
	source := FromSlice([]string{"bb", "aa", "b", "a", "cc"})
	base := source.OrderWith(func(a, b string) int { return cmp.Compare(len(a), len(b)) })
	asc := base.ThenBy(func(s string) string { return s })
	desc := base.ThenByDescending(func(s string) string { return s })
	var wg sync.WaitGroup
	for range 4 {
		wg.Go(func() {
			for _, tc := range []struct {
				q    OrderedQuery[string]
				want []string
			}{
				{base, []string{"b", "a", "bb", "aa", "cc"}},
				{asc, []string{"a", "b", "aa", "bb", "cc"}},
				{desc, []string{"b", "a", "cc", "bb", "aa"}},
			} {
				if got := tc.q.ToSlice(); !slices.Equal(got, tc.want) {
					t.Errorf("got %v want %v", got, tc.want)
				}
			}
		})
	}
	wg.Wait()
}

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

func TestOrder(t *testing.T) {
	type score int
	input := []score{3, 1, 2, 1}

	if q := Order(FromSlice(input)); !testQueryIteration(q.Query, []score{1, 1, 2, 3}) {
		t.Errorf("Order()=%v expected [1 1 2 3]", q.ToSlice())
	}
	if q := OrderDescending(FromSlice(input)); !testQueryIteration(q.Query, []score{3, 2, 1, 1}) {
		t.Errorf("OrderDescending()=%v expected [3 2 1 1]", q.ToSlice())
	}
}

func TestOrderingIsStable(t *testing.T) {
	type entry struct {
		primary   int
		secondary int
		position  int
	}

	input := make([]entry, 64)
	for i := range input {
		input[i] = entry{primary: i % 4, secondary: (i / 4) % 2, position: i}
	}

	assertStable := func(name string, q Query[entry], key func(entry) int) {
		t.Helper()
		last := make(map[int]int)
		for _, e := range q.ToSlice() {
			k := key(e)
			if previous, ok := last[k]; ok && e.position < previous {
				t.Errorf("%s reordered equal keys: position %d after %d", name, e.position, previous)
			}
			last[k] = e.position
		}
	}

	q := FromSlice(input)
	assertStable("OrderBy", q.OrderBy(func(e entry) int {
		return e.primary
	}).Query, func(e entry) int {
		return e.primary
	})
	assertStable("OrderByDescending", q.OrderByDescending(func(e entry) int {
		return e.primary
	}).Query, func(e entry) int {
		return e.primary
	})
	assertStable("ThenBy", q.OrderBy(func(e entry) int {
		return e.primary
	}).ThenBy(func(e entry) int {
		return e.secondary
	}).Query, func(e entry) int {
		return e.primary*10 + e.secondary
	})

	byPrimary := func(a, b entry) int { return cmp.Compare(a.primary, b.primary) }
	assertStable("OrderWith", q.OrderWith(byPrimary).Query, func(e entry) int {
		return e.primary
	})
	assertStable("OrderDescendingWith", q.OrderDescendingWith(byPrimary).Query, func(e entry) int {
		return e.primary
	})
}

func TestOrderWith(t *testing.T) {
	names := []string{"erin", "al", "dave", "bo", "cy", "frank"}
	byLen := func(a, b string) int { return cmp.Compare(len(a), len(b)) }

	q := FromSlice(names)
	want := []string{"al", "bo", "cy", "erin", "dave", "frank"}
	if got := q.OrderWith(byLen); !testQueryIteration(got.Query, want) {
		t.Errorf("OrderWith(byLen)=%v expected %v", got.ToSlice(), want)
	}

	wantDesc := []string{"frank", "erin", "dave", "al", "bo", "cy"}
	if got := q.OrderDescendingWith(byLen); !testQueryIteration(got.Query, wantDesc) {
		t.Errorf("OrderDescendingWith(byLen)=%v expected %v", got.ToSlice(), wantDesc)
	}

	// ThenBy orders equal-length names alphabetically instead of preserving input order.
	wantThen := []string{"al", "bo", "cy", "dave", "erin", "frank"}
	if got := q.OrderWith(byLen).ThenBy(func(s string) string { return s }); !testQueryIteration(got.Query, wantThen) {
		t.Errorf("OrderWith(byLen).ThenBy(self)=%v expected %v", got.ToSlice(), wantThen)
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
