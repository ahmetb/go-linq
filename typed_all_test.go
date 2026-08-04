package linq

import "testing"

type allNumber int

func (number allNumber) Positive() bool { return number > 0 }

func TestTypedAllAcceptsMethodExpression(t *testing.T) {
	if !fromSlice([]allNumber{1, 2, 3}).All(allNumber.Positive) {
		t.Fatal("All(allNumber.Positive) = false, want true")
	}
}

func TestTypedAllEmptyQueryIsTrue(t *testing.T) {
	if !fromSlice([]int{}).All(func(int) bool { return false }) {
		t.Fatal("All() = false for empty query, want true")
	}
}

func TestTypedAllStopsAtFirstFailure(t *testing.T) {
	calls := 0
	got := fromSlice([]int{2, 3, 4}).All(func(value int) bool {
		calls++
		return value%2 == 0
	})

	if got {
		t.Fatal("All() = true, want false")
	}
	if calls != 2 {
		t.Fatalf("predicate called %d times, want 2", calls)
	}
}

func allBenchmarkPredicate(value int) bool { return value >= 0 }

var benchmarkAllResult bool

func BenchmarkTypedAll(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkAllResult = FromSlice(source).AllT(allBenchmarkPredicate)
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkAllResult = fromSlice(source).All(allBenchmarkPredicate)
		}
	})
}
