package linq

import "testing"

func TestJoin(t *testing.T) {
	outer := []int{0, 1, 2, 3, 4, 5, 8}
	inner := []int{1, 2, 1, 4, 7, 6, 7, 2}
	want := []KeyValue[int, int]{
		{1, 1},
		{1, 1},
		{2, 2},
		{2, 2},
		{4, 4},
	}

	q := FromSlice(outer).Join(
		FromSlice(inner),
		func(i int) int { return i },
		func(i int) int { return i },
		func(outer int, inner int) KeyValue[int, int] {
			return KeyValue[int, int]{outer, inner}
		})

	if !testQueryIteration(q, want) {
		t.Errorf("FromSlice().Join()=%v expected %v", toSlice(q), want)
	}
}

func TestJoin_TypeChanging(t *testing.T) {
	type person struct {
		name   string
		petIDs []int
	}
	type pet struct {
		id   int
		name string
	}

	people := []person{{"Ahmet", []int{1, 2}}, {"Bora", []int{3}}}
	pets := []pet{{1, "Fluffy"}, {2, "Rex"}, {3, "Whiskers"}}

	q := FromSlice(people).
		SelectMany(func(p person) Query[KeyValue[string, int]] {
			return FromSlice(p.petIDs).Select(func(id int) KeyValue[string, int] {
				return KeyValue[string, int]{p.name, id}
			})
		}).
		Join(
			FromSlice(pets),
			func(kv KeyValue[string, int]) int { return kv.Value },
			func(p pet) int { return p.id },
			func(kv KeyValue[string, int], p pet) string {
				return kv.Key + ":" + p.name
			})

	want := []string{"Ahmet:Fluffy", "Ahmet:Rex", "Bora:Whiskers"}
	if !testQueryIteration(q, want) {
		t.Errorf("Join()=%v expected %v", toSlice(q), want)
	}
}

func TestLeftJoin(t *testing.T) {
	type pair struct {
		outer int
		inner int
	}

	q := FromSlice([]int{0, 1, 2, 3, 4}).LeftJoin(
		FromSlice([]int{1, 2, 1, 4, 7, 2}),
		func(i int) int { return i },
		func(i int) int { return i },
		func(outer, inner int) pair { return pair{outer, inner} },
	)
	want := []pair{{0, 0}, {1, 1}, {1, 1}, {2, 2}, {2, 2}, {3, 0}, {4, 4}}
	if !testQueryIteration(q, want) {
		t.Errorf("LeftJoin()=%v expected %v", q.ToSlice(), want)
	}
}

func TestRightJoin(t *testing.T) {
	type pair struct {
		outer int
		inner int
	}

	q := FromSlice([]int{0, 1, 2, 3, 4}).RightJoin(
		FromSlice([]int{1, 2, 1, 4, 7, 2}),
		func(i int) int { return i },
		func(i int) int { return i },
		func(outer, inner int) pair { return pair{outer, inner} },
	)
	want := []pair{{1, 1}, {2, 2}, {1, 1}, {4, 4}, {0, 7}, {2, 2}}
	if !testQueryIteration(q, want) {
		t.Errorf("RightJoin()=%v expected %v", q.ToSlice(), want)
	}
}

func TestFullJoin(t *testing.T) {
	type pair struct {
		outer int
		inner int
	}

	q := FromSlice([]int{0, 1, 2, 3, 4}).FullJoin(
		FromSlice([]int{1, 2, 1, 4, 7, 6, 7, 2}),
		func(i int) int { return i },
		func(i int) int { return i },
		func(outer, inner int) pair { return pair{outer, inner} },
	)
	want := []pair{
		{0, 0}, {1, 1}, {1, 1}, {2, 2}, {2, 2}, {3, 0}, {4, 4},
		{0, 7}, {0, 7}, {0, 6},
	}
	if !testQueryIteration(q, want) {
		t.Errorf("FullJoin()=%v expected %v", q.ToSlice(), want)
	}
}

func TestFullJoinWithEmptySide(t *testing.T) {
	type pair struct {
		outer int
		inner int
	}
	tests := []struct {
		name         string
		outer, inner []int
		want         []pair
	}{
		{name: "both", want: nil},
		{name: "inner", inner: []int{1, 2}, want: []pair{{0, 1}, {0, 2}}},
		{name: "outer", outer: []int{1, 2}, want: []pair{{1, 0}, {2, 0}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			q := FromSlice(test.outer).FullJoin(
				FromSlice(test.inner),
				func(i int) int { return i },
				func(i int) int { return i },
				func(outer, inner int) pair { return pair{outer, inner} },
			)
			if !testQueryIteration(q, test.want) {
				t.Errorf("FullJoin()=%v expected %v", q.ToSlice(), test.want)
			}
		})
	}
}

func TestFullJoinGroupsNilKeys(t *testing.T) {
	type item struct {
		key   *int
		value string
	}
	key := 1
	inner := FromSlice([]item{
		{value: "nilA"},
		{key: &key, value: "keyed"},
		{value: "nilB"},
	})

	// Nil keys never match, but unmatched inner nils share one first-seen group.
	q := FromSlice([]item{}).FullJoin(inner,
		func(item item) *int { return item.key },
		func(item item) *int { return item.key },
		func(_ item, inner item) string { return inner.value },
	)
	want := []string{"nilA", "nilB", "keyed"}
	if !testQueryIteration(q, want) {
		t.Errorf("FullJoin(mixed nil keys)=%v expected %v", q.ToSlice(), want)
	}
}

func TestJoinsIgnoreNilKeys(t *testing.T) {
	type item struct {
		key   *int
		value string
	}
	outer := FromSlice([]item{{value: "outer"}})
	inner := FromSlice([]item{{value: "inner"}})
	key := func(item item) *int { return item.key }

	if got := outer.Join(inner, key, key, func(outer, inner item) string {
		return outer.value + inner.value
	}).ToSlice(); got != nil {
		t.Errorf("Join(nil keys)=%v expected nil", got)
	}

	if q := outer.LeftJoin(inner, key, key, func(outer, inner item) string {
		return inner.value
	}); !testQueryIteration(q, []string{""}) {
		t.Errorf("LeftJoin(nil keys)=%v expected unmatched inner", q.ToSlice())
	}

	if q := outer.RightJoin(inner, key, key, func(outer, inner item) string {
		return outer.value
	}); !testQueryIteration(q, []string{""}) {
		t.Errorf("RightJoin(nil keys)=%v expected unmatched outer", q.ToSlice())
	}

	fullWant := []string{"outer:", ":inner"}
	if q := outer.FullJoin(inner, key, key, func(outer, inner item) string {
		return outer.value + ":" + inner.value
	}); !testQueryIteration(q, fullWant) {
		t.Errorf("FullJoin(nil keys)=%v expected %v", q.ToSlice(), fullWant)
	}
	if q := outer.FullJoin(
		inner,
		func(item item) any { return item.key },
		func(item item) any { return item.key },
		func(outer, inner item) string { return outer.value + ":" + inner.value },
	); !testQueryIteration(q, fullWant) {
		t.Errorf("FullJoin(typed nil interface keys)=%v expected %v", q.ToSlice(), fullWant)
	}
	if q := outer.FullJoin(
		inner,
		func(item) any { return nil },
		func(item) any { return nil },
		func(outer, inner item) string { return outer.value + ":" + inner.value },
	); !testQueryIteration(q, fullWant) {
		t.Errorf("FullJoin(untyped nil interface keys)=%v expected %v", q.ToSlice(), fullWant)
	}

	if q := outer.GroupJoin(inner, key, key, func(_ item, inner []item) int {
		return len(inner)
	}); !testQueryIteration(q, []int{0}) {
		t.Errorf("GroupJoin(nil keys)=%v expected [0]", q.ToSlice())
	}

	interfaceLookup := buildJoinLookup(inner, func(item item) any { return item.key })
	if len(interfaceLookup) != 0 {
		t.Errorf("join lookup with a typed nil interface key=%v expected empty", interfaceLookup)
	}

	type boxedItem struct {
		key   any
		value string
	}
	boxedOuter := FromSlice([]boxedItem{{key: []int(nil), value: "outer"}})
	boxedInner := FromSlice([]boxedItem{{key: map[string]int(nil), value: "inner"}})
	if q := boxedOuter.FullJoin(
		boxedInner,
		func(item boxedItem) any { return item.key },
		func(item boxedItem) any { return item.key },
		func(outer, inner boxedItem) string { return outer.value + ":" + inner.value },
	); !testQueryIteration(q, fullWant) {
		t.Errorf("FullJoin(non-comparable typed nil interface keys)=%v expected %v", q.ToSlice(), fullWant)
	}
}

func TestJoinsIgnoreNonHashableNilKeys(t *testing.T) {
	type item struct {
		key   any
		value string
	}
	key := func(item item) any { return item.key }
	// The searched side must hold a real key: on an empty lookup Go answers
	// before it hashes, which hides an unhashable key.
	keyed := FromSlice([]item{{key: 1, value: "keyed"}})
	boxedNil := FromSlice([]item{{key: []int(nil), value: "nil"}})
	pair := func(outer, inner item) string { return outer.value + ":" + inner.value }

	if got := boxedNil.Join(keyed, key, key, pair).ToSlice(); got != nil {
		t.Errorf("Join(boxed nil key)=%v expected nil", got)
	}
	if q := boxedNil.LeftJoin(keyed, key, key, pair); !testQueryIteration(q, []string{"nil:"}) {
		t.Errorf("LeftJoin(boxed nil key)=%v expected [nil:]", q.ToSlice())
	}
	if q := keyed.RightJoin(boxedNil, key, key, pair); !testQueryIteration(q, []string{":nil"}) {
		t.Errorf("RightJoin(boxed nil key)=%v expected [:nil]", q.ToSlice())
	}
	if q := boxedNil.GroupJoin(keyed, key, key, func(_ item, inner []item) int {
		return len(inner)
	}); !testQueryIteration(q, []int{0}) {
		t.Errorf("GroupJoin(boxed nil key)=%v expected [0]", q.ToSlice())
	}
	want := []string{"nil:", ":keyed"}
	if q := boxedNil.FullJoin(keyed, key, key, pair); !testQueryIteration(q, want) {
		t.Errorf("FullJoin(boxed nil key)=%v expected %v", q.ToSlice(), want)
	}
}

func TestOuterJoinDoesNotReadOtherSideWhenRetainedSideIsEmpty(t *testing.T) {
	leftPulled, rightPulled := 0, 0
	left := FromSlice([]int{1}).Where(func(int) bool {
		leftPulled++
		return true
	})
	right := FromSlice([]int{1}).Where(func(int) bool {
		rightPulled++
		return true
	})

	FromSlice([]int{}).LeftJoin(
		right,
		func(i int) int { return i },
		func(i int) int { return i },
		func(outer, inner int) int { return outer + inner },
	).ToSlice()
	left.RightJoin(
		FromSlice([]int{}),
		func(i int) int { return i },
		func(i int) int { return i },
		func(outer, inner int) int { return outer + inner },
	).ToSlice()

	if rightPulled != 0 || leftPulled != 0 {
		t.Errorf("empty retained side pulled left=%d right=%d expected 0,0", leftPulled, rightPulled)
	}
}

func TestLeftJoinStopsWithConsumer(t *testing.T) {
	outerPulled, innerPulled := 0, 0
	outer := FromSlice([]int{1, 2, 3}).Where(func(int) bool {
		outerPulled++
		return true
	})
	inner := FromSlice([]int{1, 2, 3}).Where(func(int) bool {
		innerPulled++
		return true
	})

	outer.LeftJoin(
		inner,
		func(i int) int { return i },
		func(i int) int { return i },
		func(outer, inner int) int { return outer + inner },
	).Iterate(func(int) bool { return false })

	if outerPulled != 1 || innerPulled != 3 {
		t.Errorf("LeftJoin pulled outer=%d inner=%d expected 1,3", outerPulled, innerPulled)
	}
}

func TestFullJoinStopsWithConsumer(t *testing.T) {
	outerPulled, innerPulled, results := 0, 0, 0
	outer := FromSlice([]int{1, 2, 3}).Where(func(int) bool {
		outerPulled++
		return true
	})
	inner := FromSlice([]int{1, 2, 3}).Where(func(int) bool {
		innerPulled++
		return true
	})

	outer.FullJoin(
		inner,
		func(i int) int { return i },
		func(i int) int { return i },
		func(outer, inner int) int {
			results++
			return outer + inner
		},
	).Iterate(func(int) bool { return false })

	if outerPulled != 1 || innerPulled != 3 || results != 1 {
		t.Errorf("FullJoin pulled outer=%d inner=%d results=%d expected 1,3,1",
			outerPulled, innerPulled, results)
	}
}

func TestLeftJoinWithEmptyInner(t *testing.T) {
	outerPulled := 0
	outer := FromSlice([]int{1, 2, 3}).Where(func(int) bool {
		outerPulled++
		return true
	})

	q := outer.LeftJoin(
		FromSlice([]int{}),
		func(i int) int { return i },
		func(i int) int { return i },
		func(outer, inner int) KeyValue[int, int] {
			return KeyValue[int, int]{outer, inner}
		},
	)

	want := []KeyValue[int, int]{{1, 0}, {2, 0}, {3, 0}}
	if !testQueryIteration(q, want) {
		t.Errorf("LeftJoin(empty inner)=%v expected %v", q.ToSlice(), want)
	}
	if outerPulled == 0 {
		t.Error("LeftJoin(empty inner) skipped outer; every outer element owes a result")
	}
}

func TestRightJoinWithEmptyOuter(t *testing.T) {
	q := FromSlice([]int{}).RightJoin(
		FromSlice([]int{1, 2, 3}),
		func(i int) int { return i },
		func(i int) int { return i },
		func(outer, inner int) KeyValue[int, int] {
			return KeyValue[int, int]{outer, inner}
		},
	)

	want := []KeyValue[int, int]{{0, 1}, {0, 2}, {0, 3}}
	if !testQueryIteration(q, want) {
		t.Errorf("RightJoin(empty outer)=%v expected %v", q.ToSlice(), want)
	}
}
