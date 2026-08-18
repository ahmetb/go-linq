package linq

import "testing"

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

func TestDerivedSizeHints(t *testing.T) {
	q := FromSlice(make([]int, 10))
	checkSizeHint(t, "SkipLast", q.SkipLast(3), 7)
	checkSizeHint(t, "TakeLast", q.TakeLast(3), 3)
	checkSizeHint(t, "TakeRange", q.TakeRange(IndexFromEnd(7), IndexFromStart(8)), 5)
	checkSizeHint(t, "Zip3", q.Zip3(Range(0, 8), Range(0, 6), func(a, b, c int) int {
		return a + b + c
	}), 6)
	checkSizeHint(t, "Chunk", Chunk(q, 3), 4)
}
