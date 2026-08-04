package linq

import "testing"

func TestTypedCountWithAcceptsMethodExpression(t *testing.T) {
	got := fromSlice([]allNumber{-1, 2, 3, -4}).CountWith(allNumber.Positive)
	if got != 2 {
		t.Fatalf("CountWith(allNumber.Positive) = %d, want 2", got)
	}
}

func TestTypedCountWithEvaluatesEveryElement(t *testing.T) {
	calls := 0
	got := fromSlice([]int{1, 2, 3}).CountWith(func(value int) bool {
		calls++
		return value%2 == 1
	})

	if got != 2 {
		t.Fatalf("CountWith() = %d, want 2", got)
	}
	if calls != 3 {
		t.Fatalf("predicate called %d times, want 3", calls)
	}
}

func countWithBenchmarkPredicate(value int) bool { return value%2 == 0 }

var benchmarkCountWithResult int

func BenchmarkTypedCountWith(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkCountWithResult = FromSlice(source).CountWithT(countWithBenchmarkPredicate)
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkCountWithResult = fromSlice(source).CountWith(countWithBenchmarkPredicate)
		}
	})
}
