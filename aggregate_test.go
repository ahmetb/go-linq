package linq

import "testing"

func aggregateSum(left, right int) int { return left + right }

func TestTypedAggregate(t *testing.T) {
	if got := fromSlice([]int{1, 2, 3, 4}).Aggregate(aggregateSum); got != 10 {
		t.Fatalf("Aggregate() = %d, want 10", got)
	}
}

func TestTypedAggregateEmptyReturnsZero(t *testing.T) {
	if got := fromSlice([]int(nil)).Aggregate(aggregateSum); got != 0 {
		t.Fatalf("empty Aggregate() = %d, want 0", got)
	}
}

var benchmarkAggregateResult int

func BenchmarkTypedAggregate(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkAggregateResult = legacyFromSlice(source).AggregateT(aggregateSum).(int)
		}
	})

	b.Run("generic_query", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkAggregateResult = fromSlice(source).Aggregate(aggregateSum)
		}
	})
}
