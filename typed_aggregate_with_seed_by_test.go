package linq

import (
	"strconv"
	"testing"
)

func formatAggregate(value int) string { return "sum:" + strconv.Itoa(value) }

func TestTypedAggregateWithSeedByDifferentResultType(t *testing.T) {
	got := fromSlice([]int{1, 2, 3}).AggregateWithSeedBy(0, aggregateSum, formatAggregate)
	if got != "sum:6" {
		t.Fatalf("AggregateWithSeedBy() = %q, want sum:6", got)
	}
}

func TestTypedAggregateWithSeedByEmptySelectsSeed(t *testing.T) {
	got := fromSlice([]int(nil)).AggregateWithSeedBy(10, aggregateSum, formatAggregate)
	if got != "sum:10" {
		t.Fatalf("empty AggregateWithSeedBy() = %q, want sum:10", got)
	}
}

func aggregateIdentity(value int) int { return value }

var benchmarkAggregateWithSeedByResult int

func BenchmarkTypedAggregateWithSeedBy(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkAggregateWithSeedByResult = legacyFromSlice(source).
				AggregateWithSeedByT(0, aggregateSum, aggregateIdentity).(int)
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkAggregateWithSeedByResult = fromSlice(source).
				AggregateWithSeedBy(0, aggregateSum, aggregateIdentity)
		}
	})
}
