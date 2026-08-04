package linq

import "testing"

func TestTypedIndexOfAcceptsMethodExpression(t *testing.T) {
	got := fromSlice([]allNumber{-1, -2, 3, 4}).IndexOf(allNumber.Positive)
	if got != 2 {
		t.Fatalf("IndexOf(allNumber.Positive) = %d, want 2", got)
	}
}

func TestTypedIndexOfReturnsNegativeOneWithoutMatch(t *testing.T) {
	got := fromSlice([]int{1, 3, 5}).IndexOf(func(value int) bool { return value%2 == 0 })
	if got != -1 {
		t.Fatalf("IndexOf() = %d, want -1", got)
	}
}

func TestTypedIndexOfStopsAtFirstMatch(t *testing.T) {
	calls := 0
	got := fromSlice([]int{1, 2, 4}).IndexOf(func(value int) bool {
		calls++
		return value%2 == 0
	})

	if got != 1 {
		t.Fatalf("IndexOf() = %d, want 1", got)
	}
	if calls != 2 {
		t.Fatalf("predicate called %d times, want 2", calls)
	}
}

func indexOfBenchmarkPredicate(value int) bool { return value == 1023 }

var benchmarkIndexOfResult int

func BenchmarkTypedIndexOf(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkIndexOfResult = FromSlice(source).IndexOfT(indexOfBenchmarkPredicate)
		}
	})

	b.Run("generic_query", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkIndexOfResult = fromSlice(source).IndexOf(indexOfBenchmarkPredicate)
		}
	})
}
