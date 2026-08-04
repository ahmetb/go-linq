package linq

import (
	"slices"
	"testing"
)

func TestTypedIterateExposesTypedSequence(t *testing.T) {
	var got []int
	for value := range fromSlice([]int{1, 2, 3}).Iterate {
		got = append(got, value)
	}

	if !slices.Equal(got, []int{1, 2, 3}) {
		t.Fatalf("Iterate produced %v, want [1 2 3]", got)
	}
}

func TestTypedIteratePropagatesConsumerStop(t *testing.T) {
	calls := 0
	fromSlice([]string{"one", "two"}).Iterate(func(string) bool {
		calls++
		return false
	})

	if calls != 1 {
		t.Fatalf("consumer called %d times, want 1", calls)
	}
}

var benchmarkIterateSum int

func BenchmarkTypedIterate(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			legacyFromSlice(source).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkIterateSum = sum
		}
	})

	b.Run("generic_query", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).Iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkIterateSum = sum
		}
	})
}
