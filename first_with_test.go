package linq

import "testing"

func TestTypedFirstWithAcceptsMethodExpression(t *testing.T) {
	got := fromSlice([]allNumber{-1, 2, 3}).FirstWith(allNumber.Positive)
	if got != 2 {
		t.Fatalf("FirstWith(allNumber.Positive) = %d, want 2", got)
	}
}

func TestTypedFirstWithReturnsZeroWithoutMatch(t *testing.T) {
	got := fromSlice([]string{"one", "two"}).FirstWith(func(value string) bool { return len(value) > 10 })
	if got != "" {
		t.Fatalf("FirstWith() = %q, want zero value", got)
	}
}

func TestTypedFirstWithStopsAtFirstMatch(t *testing.T) {
	calls := 0
	got := fromSlice([]int{1, 2, 4}).FirstWith(func(value int) bool {
		calls++
		return value%2 == 0
	})

	if got != 2 {
		t.Fatalf("FirstWith() = %d, want 2", got)
	}
	if calls != 2 {
		t.Fatalf("predicate called %d times, want 2", calls)
	}
}

func firstWithBenchmarkPredicate(value int) bool { return value == 1023 }

var benchmarkFirstWithResult int

func BenchmarkTypedFirstWith(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkFirstWithResult = legacyFromSlice(source).FirstWithT(firstWithBenchmarkPredicate).(int)
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkFirstWithResult = fromSlice(source).FirstWith(firstWithBenchmarkPredicate)
		}
	})
}
