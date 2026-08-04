package linq

import (
	"slices"
	"testing"
)

func TestTypedSkipBypassesRequestedCount(t *testing.T) {
	got := fromSlice([]int{1, 2, 3, 4}).Skip(2).toSlice()
	if !slices.Equal(got, []int{3, 4}) {
		t.Fatalf("Skip(2) = %v, want [3 4]", got)
	}
}

func TestTypedSkipNonPositiveCountReturnsAllValues(t *testing.T) {
	got := fromSlice([]int{1, 2}).Skip(-1).toSlice()
	if !slices.Equal(got, []int{1, 2}) {
		t.Fatalf("Skip(-1) = %v, want [1 2]", got)
	}
}

func TestTypedSkipStopsWithConsumer(t *testing.T) {
	var got []int
	fromSlice([]int{1, 2, 3}).Skip(1).iterate(func(value int) bool {
		got = append(got, value)
		return false
	})

	if !slices.Equal(got, []int{2}) {
		t.Fatalf("Skip(1) yielded %v, want [2]", got)
	}
}

var benchmarkSkipSum int

func BenchmarkTypedSkip(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			legacyFromSlice(source).Skip(512).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkSkipSum = sum
		}
	})

	b.Run("generic_query", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).Skip(512).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkSkipSum = sum
		}
	})
}
