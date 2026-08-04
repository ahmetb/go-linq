package linq

import "testing"

func TestTypedAnyWithAcceptsMethodExpression(t *testing.T) {
	if !fromSlice([]allNumber{-1, 2, 3}).AnyWith(allNumber.Positive) {
		t.Fatal("AnyWith(allNumber.Positive) = false, want true")
	}
}

func TestTypedAnyWithReturnsFalseWithoutMatch(t *testing.T) {
	if fromSlice([]int{1, 3, 5}).AnyWith(func(value int) bool { return value%2 == 0 }) {
		t.Fatal("AnyWith() = true, want false")
	}
}

func TestTypedAnyWithStopsAtFirstMatch(t *testing.T) {
	calls := 0
	got := fromSlice([]int{1, 2, 4}).AnyWith(func(value int) bool {
		calls++
		return value%2 == 0
	})

	if !got {
		t.Fatal("AnyWith() = false, want true")
	}
	if calls != 2 {
		t.Fatalf("predicate called %d times, want 2", calls)
	}
}

func anyWithBenchmarkPredicate(value int) bool { return value < 0 }

var benchmarkAnyWithResult bool

func BenchmarkTypedAnyWith(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkAnyWithResult = legacyFromSlice(source).AnyWithT(anyWithBenchmarkPredicate)
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkAnyWithResult = fromSlice(source).AnyWith(anyWithBenchmarkPredicate)
		}
	})
}
