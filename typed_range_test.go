package linq

import (
	"slices"
	"testing"
)

func TestTypedRangeGeneratesRequestedValues(t *testing.T) {
	got := integerRange(3, 4).toSlice()
	if !slices.Equal(got, []int{3, 4, 5, 6}) {
		t.Fatalf("integerRange(3, 4) = %v, want [3 4 5 6]", got)
	}
}

func TestTypedRangeNonPositiveCountIsEmpty(t *testing.T) {
	if got := integerRange(3, -1).toSlice(); len(got) != 0 {
		t.Fatalf("integerRange(3, -1) = %v, want empty", got)
	}
}

func TestTypedRangeStopsWithConsumer(t *testing.T) {
	var got []int
	integerRange(10, 3).iterate(func(value int) bool {
		got = append(got, value)
		return false
	})

	if !slices.Equal(got, []int{10}) {
		t.Fatalf("integerRange() yielded %v, want [10]", got)
	}
}

var benchmarkRangeSum int

func BenchmarkTypedRange(b *testing.B) {
	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			Range(0, 1024).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkRangeSum = sum
		}
	})

	b.Run("generic_query", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			integerRange(0, 1024).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkRangeSum = sum
		}
	})
}
