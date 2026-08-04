package linq

import (
	"strconv"
	"testing"
)

func appendNumber(accumulator string, value int) string {
	return accumulator + strconv.Itoa(value)
}

func TestTypedAggregateWithSeedDifferentAccumulatorType(t *testing.T) {
	if got := fromSlice([]int{1, 2, 3}).AggregateWithSeed("values:", appendNumber); got != "values:123" {
		t.Fatalf("AggregateWithSeed() = %q, want values:123", got)
	}
}

func TestTypedAggregateWithSeedEmptyReturnsSeed(t *testing.T) {
	if got := fromSlice([]int(nil)).AggregateWithSeed("seed", appendNumber); got != "seed" {
		t.Fatalf("empty AggregateWithSeed() = %q, want seed", got)
	}
}

var benchmarkAggregateWithSeedResult int

func BenchmarkTypedAggregateWithSeed(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkAggregateWithSeedResult = FromSlice(source).AggregateWithSeedT(0, aggregateSum).(int)
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkAggregateWithSeedResult = fromSlice(source).AggregateWithSeed(0, aggregateSum)
		}
	})
}
