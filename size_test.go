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

// TestSizeHint_ComputedByCountDeterministicOperators pins the operators whose
// output count is a deterministic function of their inputs' counts: they must
// compute the hint (cap == len after collecting, see
// TestSizeHint_ExactPreallocation for why that detects preallocation).
func TestSizeHint_ComputedByCountDeterministicOperators(t *testing.T) {
	src := func() Query[int] { return FromSlice(make([]int, 1000)) }
	tests := []struct {
		name string
		q    Query[int]
		want int
	}{
		{"Take", src().Take(600), 600},
		{"TakeBeyondLen", src().Take(2000), 1000},
		{"Skip", src().Skip(400), 600},
		{"Concat", src().Concat(src()), 2000},
		{"Append", src().Append(1), 1001},
		{"Prepend", src().Prepend(1), 1001},
		{"Zip", src().Zip(src().Take(500), func(a, b int) int { return a + b }), 500},
		{"DefaultIfEmpty", src().DefaultIfEmpty(0), 1000},
	}
	for _, test := range tests {
		out := test.q.ToSlice()
		if len(out) != test.want {
			t.Fatalf("%s: len=%d expected %d", test.name, len(out), test.want)
		}
		if cap(out) != len(out) {
			t.Errorf("%s: cap=%d len=%d: hint was not computed; result grew instead of preallocating",
				test.name, cap(out), len(out))
		}
	}
}

// TestSizeHint_TakeOfUnknownSourceDoesNotOverallocate catches Take using its
// count alone as a hint: count is an upper bound, not a yield count, and must
// not preallocate when the source size is unknown.
func TestSizeHint_TakeOfUnknownSourceDoesNotOverallocate(t *testing.T) {
	out := FromSlice(make([]int, 10)).Where(func(int) bool {
		return true // keep everything, but the count is no longer statically known
	}).Take(100_000).ToSlice()

	if len(out) != 10 {
		t.Fatalf("len=%d expected 10", len(out))
	}
	if cap(out) >= 100_000 {
		t.Errorf("cap=%d: Take preallocated its count with an unknown source size", cap(out))
	}
}
