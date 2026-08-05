package linq

import (
	"testing"
)

// The set operators pick a strongly-typed set for basic element kinds and
// fall back to boxed comparison otherwise. These tests exercise the fallback
// paths and the panic behavior for non-comparable element types.

func TestDistinct_StructFallback(t *testing.T) {
	type point struct{ x, y int }

	input := []point{{1, 1}, {2, 2}, {1, 1}, {3, 3}, {2, 2}}
	want := []point{{1, 1}, {2, 2}, {3, 3}}

	if q := FromSlice(input).Distinct(); !testQueryIteration(q, want) {
		t.Errorf("Distinct()=%v expected %v", toSlice(q), want)
	}
}

func TestDistinct_NamedTypeFallback(t *testing.T) {
	type id int

	input := []id{1, 2, 2, 3, 1}
	want := []id{1, 2, 3}

	if q := FromSlice(input).Distinct(); !testQueryIteration(q, want) {
		t.Errorf("Distinct()=%v expected %v", toSlice(q), want)
	}
}

func TestDistinct_NonComparablePanics(t *testing.T) {
	mustPanic(t, func() {
		FromSlice([][]int{{1}, {2}}).Distinct().Count()
	})
}

func TestSetOps_StructFallback(t *testing.T) {
	type point struct{ x, y int }

	a := []point{{1, 1}, {2, 2}, {3, 3}}
	b := []point{{2, 2}, {4, 4}}

	wantUnion := []point{{1, 1}, {2, 2}, {3, 3}, {4, 4}}
	if q := FromSlice(a).Union(FromSlice(b)); !testQueryIteration(q, wantUnion) {
		t.Errorf("Union()=%v expected %v", toSlice(q), wantUnion)
	}

	wantExcept := []point{{1, 1}, {3, 3}}
	if q := FromSlice(a).Except(FromSlice(b)); !testQueryIteration(q, wantExcept) {
		t.Errorf("Except()=%v expected %v", toSlice(q), wantExcept)
	}

	wantIntersect := []point{{2, 2}}
	if q := FromSlice(a).Intersect(FromSlice(b)); !testQueryIteration(q, wantIntersect) {
		t.Errorf("Intersect()=%v expected %v", toSlice(q), wantIntersect)
	}
}

func TestContains_StructFallback(t *testing.T) {
	type point struct{ x, y int }

	q := FromSlice([]point{{1, 1}, {2, 2}})
	if !q.Contains(point{2, 2}) {
		t.Error("Contains({2,2})=false expected true")
	}
	if q.Contains(point{3, 3}) {
		t.Error("Contains({3,3})=true expected false")
	}
}

func TestSequenceEqual_StructFallback(t *testing.T) {
	type point struct{ x, y int }

	a := FromSlice([]point{{1, 1}, {2, 2}})
	if !a.SequenceEqual(FromSlice([]point{{1, 1}, {2, 2}})) {
		t.Error("SequenceEqual()=false expected true")
	}
	if a.SequenceEqual(FromSlice([]point{{1, 1}, {2, 3}})) {
		t.Error("SequenceEqual()=true expected false")
	}
}

// TestSort_SingleComparatorCall verifies the less function is invoked exactly
// once per comparison (sorting two elements requires exactly one comparison).
func TestSort_SingleComparatorCall(t *testing.T) {
	calls := 0
	FromSlice([]int{2, 1}).Sort(func(i, j int) bool {
		calls++
		return i < j
	}).Count()

	if calls != 1 {
		t.Errorf("less called %d times for a 2-element sort, expected 1", calls)
	}
}
