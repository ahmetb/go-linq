package linq

import "testing"

// The size hint carried by Query is an invariant the compiler cannot check:
// it may be propagated only by operators that emit exactly one element per
// source element. These tests pin it in both directions.

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
