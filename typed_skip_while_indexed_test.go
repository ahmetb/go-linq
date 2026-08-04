package linq

import (
	"slices"
	"testing"
)

func TestTypedSkipWhileIndexedSuppliesSourceIndex(t *testing.T) {
	var indices []int
	got := fromSlice([]string{"a", "b", "c", "d"}).SkipWhileIndexed(func(index int, _ string) bool {
		indices = append(indices, index)
		return index < 2
	}).toSlice()

	if !slices.Equal(got, []string{"c", "d"}) {
		t.Fatalf("SkipWhileIndexed() = %v, want [c d]", got)
	}
	if !slices.Equal(indices, []int{0, 1, 2}) {
		t.Fatalf("predicate indices = %v, want [0 1 2]", indices)
	}
}

func TestTypedSkipWhileIndexedStopsWithConsumer(t *testing.T) {
	calls := 0
	q := fromSlice([]int{1, 2, 3}).SkipWhileIndexed(func(index, _ int) bool {
		calls++
		return index == 0
	})

	q.iterate(func(int) bool { return false })

	if calls != 2 {
		t.Fatalf("predicate called %d times, want 2", calls)
	}
}

func skipWhileIndexedBenchmarkPredicate(index, _ int) bool { return index < 512 }

var benchmarkSkipWhileIndexedSum int

func BenchmarkTypedSkipWhileIndexed(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			legacyFromSlice(source).SkipWhileIndexedT(skipWhileIndexedBenchmarkPredicate).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkSkipWhileIndexedSum = sum
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).SkipWhileIndexed(skipWhileIndexedBenchmarkPredicate).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkSkipWhileIndexedSum = sum
		}
	})
}
