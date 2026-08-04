package linq

import (
	"slices"
	"testing"
)

func TestTypedTakeWhileIndexedSuppliesSourceIndex(t *testing.T) {
	var indices []int
	got := fromSlice([]string{"a", "b", "c", "d"}).TakeWhileIndexed(func(index int, _ string) bool {
		indices = append(indices, index)
		return index < 2
	}).toSlice()

	if !slices.Equal(got, []string{"a", "b"}) {
		t.Fatalf("TakeWhileIndexed() = %v, want [a b]", got)
	}
	if !slices.Equal(indices, []int{0, 1, 2}) {
		t.Fatalf("predicate indices = %v, want [0 1 2]", indices)
	}
}

func TestTypedTakeWhileIndexedStopsWithConsumer(t *testing.T) {
	calls := 0
	q := fromSlice([]int{1, 2, 3}).TakeWhileIndexed(func(int, int) bool {
		calls++
		return true
	})

	q.iterate(func(int) bool { return false })

	if calls != 1 {
		t.Fatalf("predicate called %d times, want 1", calls)
	}
}

func takeWhileIndexedBenchmarkPredicate(index, _ int) bool { return index < 512 }

var benchmarkTakeWhileIndexedSum int

func BenchmarkTypedTakeWhileIndexed(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			legacyFromSlice(source).TakeWhileIndexedT(takeWhileIndexedBenchmarkPredicate).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkTakeWhileIndexedSum = sum
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).TakeWhileIndexed(takeWhileIndexedBenchmarkPredicate).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkTakeWhileIndexedSum = sum
		}
	})
}
