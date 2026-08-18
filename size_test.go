package linq

import (
	"math"
	"slices"
	"testing"
)

// The size hint carried by Query is an invariant the compiler cannot check:
// it may be propagated only when an operator can derive its output count
// exactly. These tests pin it in both directions.

// TestSizeHint_ExactPreallocation catches an operator losing the hint: append
// growth never lands on an arbitrary exact size (collecting 1000 elements by
// doubling ends at capacity 1024), so cap == len is only achievable when the
// result was preallocated from the hint.
func TestSizeHint_ExactPreallocation(t *testing.T) {
	out := FromSlice(make([]int, 1000)).Select(func(i int) int {
		return i * 2
	}).ToSlice()

	if len(out) != 1000 {
		t.Fatalf("len=%d expected 1000", len(out))
	}
	if cap(out) != len(out) {
		t.Errorf("cap=%d len=%d: size hint was lost; result grew instead of preallocating",
			cap(out), len(out))
	}
}

// TestSizeHint_DroppedByFilters catches an operator wrongly inheriting the
// hint: a filtered query must not preallocate source-sized results.
func TestSizeHint_DroppedByFilters(t *testing.T) {
	out := FromSlice(make([]int, 100_000)).Where(func(i int) bool {
		return false // keep nothing
	}).Append(1).ToSlice()

	if len(out) != 1 {
		t.Fatalf("len=%d expected 1", len(out))
	}
	if cap(out) >= 100_000 {
		t.Errorf("cap=%d: filtered query inherited the source's size hint", cap(out))
	}
}

func checkSizeHint[T any](t *testing.T, name string, q Query[T], want int) {
	t.Helper()
	if got := q.size; got != want {
		t.Errorf("%s size=%d expected %d", name, got, want)
	}
}

// TestDerivedSizeHints pins the operators whose output count follows exactly
// from their input sizes.
func TestDerivedSizeHints(t *testing.T) {
	q := FromSlice(make([]int, 10))
	checkSizeHint(t, "Take", q.Take(3), 3)
	checkSizeHint(t, "Take past end", q.Take(50), 10)
	checkSizeHint(t, "Skip", q.Skip(3), 7)
	checkSizeHint(t, "Skip past end", q.Skip(50), 0)
	checkSizeHint(t, "Skip negative", q.Skip(-5), 10)
	checkSizeHint(t, "Concat", q.Concat(Range(0, 5)), 15)
	checkSizeHint(t, "Append", q.Append(1), 11)
	checkSizeHint(t, "Prepend", q.Prepend(1), 11)
	checkSizeHint(t, "DefaultIfEmpty", q.DefaultIfEmpty(0), 10)
	checkSizeHint(t, "Zip", q.Zip(Range(0, 4), func(a, b int) int {
		return a + b
	}), 4)
}

// TestUnknownSizeHints pins the other direction: no operator may invent a size
// for an input that carries none.
func TestUnknownSizeHints(t *testing.T) {
	sized := FromSlice(make([]int, 10))
	unsized := FromSeq(slices.Values(make([]int, 10)))

	checkSizeHint(t, "unsized source", unsized, 0)
	checkSizeHint(t, "Concat unsized right", sized.Concat(unsized), 0)
	checkSizeHint(t, "Concat unsized left", unsized.Concat(sized), 0)
	checkSizeHint(t, "Append unsized", unsized.Append(1), 0)
	checkSizeHint(t, "Append overflow", Repeat(0, math.MaxInt).Append(0), 0)
	checkSizeHint(t, "Prepend unsized", unsized.Prepend(1), 0)
	checkSizeHint(t, "Take unsized", unsized.Take(3), 0)
	checkSizeHint(t, "Skip unsized", unsized.Skip(3), 0)
	checkSizeHint(t, "DefaultIfEmpty unsized", unsized.DefaultIfEmpty(0), 0)
	checkSizeHint(t, "Zip unsized", sized.Zip(unsized, func(a, b int) int {
		return a + b
	}), 0)
	checkSizeHint(t, "Where", sized.Where(func(int) bool {
		return true
	}), 0)
}
