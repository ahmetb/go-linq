package linq

import "testing"

func TestTypedLastWithAcceptsMethodExpression(t *testing.T) {
	got := fromSlice([]allNumber{-1, 2, 3, -4}).LastWith(allNumber.Positive)
	if got != 3 {
		t.Fatalf("LastWith(allNumber.Positive) = %d, want 3", got)
	}
}

func TestTypedLastWithReturnsZeroWithoutMatch(t *testing.T) {
	got := fromSlice([]string{"one", "two"}).LastWith(func(value string) bool { return len(value) > 10 })
	if got != "" {
		t.Fatalf("LastWith() = %q, want zero value", got)
	}
}

func TestTypedLastWithEvaluatesEveryElement(t *testing.T) {
	calls := 0
	got := fromSlice([]int{2, 3, 4}).LastWith(func(value int) bool {
		calls++
		return value%2 == 0
	})

	if got != 4 {
		t.Fatalf("LastWith() = %d, want 4", got)
	}
	if calls != 3 {
		t.Fatalf("predicate called %d times, want 3", calls)
	}
}

func lastWithBenchmarkPredicate(value int) bool { return value%2 == 0 }

var benchmarkLastWithResult int

func BenchmarkTypedLastWith(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkLastWithResult = legacyFromSlice(source).LastWithT(lastWithBenchmarkPredicate).(int)
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkLastWithResult = fromSlice(source).LastWith(lastWithBenchmarkPredicate)
		}
	})
}
