package linq

import "testing"

func TestTypedSingleWithAcceptsMethodExpression(t *testing.T) {
	got := fromSlice([]allNumber{-1, 2, -3}).SingleWith(allNumber.Positive)
	if got != 2 {
		t.Fatalf("SingleWith(allNumber.Positive) = %d, want 2", got)
	}
}

func TestTypedSingleWithReturnsZeroForMultipleMatches(t *testing.T) {
	calls := 0
	got := fromSlice([]int{2, 3, 4, 6}).SingleWith(func(value int) bool {
		calls++
		return value%2 == 0
	})

	if got != 0 {
		t.Fatalf("SingleWith() = %d, want zero value", got)
	}
	if calls != 3 {
		t.Fatalf("predicate called %d times, want 3 (through the second match)", calls)
	}
}

func TestTypedSingleWithReturnsZeroWithoutMatch(t *testing.T) {
	got := fromSlice([]string{"one", "two"}).SingleWith(func(value string) bool { return len(value) > 10 })
	if got != "" {
		t.Fatalf("SingleWith() = %q, want zero value", got)
	}
}

func singleWithBenchmarkPredicate(value int) bool { return value == 1023 }

var benchmarkSingleWithResult int

func BenchmarkTypedSingleWith(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkSingleWithResult = legacyFromSlice(source).SingleWithT(singleWithBenchmarkPredicate).(int)
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkSingleWithResult = fromSlice(source).SingleWith(singleWithBenchmarkPredicate)
		}
	})
}
