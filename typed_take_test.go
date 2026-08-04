package linq

import (
	"slices"
	"testing"
)

func TestTypedTakeYieldsRequestedPrefix(t *testing.T) {
	got := fromSlice([]int{1, 2, 3, 4}).Take(2).toSlice()
	if !slices.Equal(got, []int{1, 2}) {
		t.Fatalf("Take(2) = %v, want [1 2]", got)
	}
}

func TestTypedTakeNonPositiveCountDoesNotEnumerateSource(t *testing.T) {
	visited := 0
	source := Query[int]{iterate: func(func(int) bool) { visited++ }}

	if got := source.Take(0).toSlice(); len(got) != 0 {
		t.Fatalf("Take(0) = %v, want empty", got)
	}
	if visited != 0 {
		t.Fatalf("source enumerated %d times, want 0", visited)
	}
}

func TestTypedTakeStopsWithConsumer(t *testing.T) {
	var got []int
	fromSlice([]int{1, 2, 3}).Take(2).iterate(func(value int) bool {
		got = append(got, value)
		return false
	})

	if !slices.Equal(got, []int{1}) {
		t.Fatalf("Take(2) yielded %v, want [1]", got)
	}
}

var benchmarkTakeSum int

func BenchmarkTypedTake(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			FromSlice(source).Take(512).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkTakeSum = sum
		}
	})

	b.Run("generic_query", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).Take(512).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkTakeSum = sum
		}
	})
}
