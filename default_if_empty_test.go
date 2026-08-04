package linq

import (
	"slices"
	"testing"
)

func TestTypedDefaultIfEmptyYieldsDefaultForEmptySource(t *testing.T) {
	got := fromSlice([]string{}).DefaultIfEmpty("default").toSlice()
	if !slices.Equal(got, []string{"default"}) {
		t.Fatalf("DefaultIfEmpty() = %v, want [default]", got)
	}
}

func TestTypedDefaultIfEmptyPreservesNonEmptySource(t *testing.T) {
	got := fromSlice([]int{1, 2}).DefaultIfEmpty(9).toSlice()
	if !slices.Equal(got, []int{1, 2}) {
		t.Fatalf("DefaultIfEmpty(9) = %v, want [1 2]", got)
	}
}

func TestTypedDefaultIfEmptyStopsWithConsumer(t *testing.T) {
	var got []int
	fromSlice([]int{1, 2}).DefaultIfEmpty(9).iterate(func(value int) bool {
		got = append(got, value)
		return false
	})

	if !slices.Equal(got, []int{1}) {
		t.Fatalf("DefaultIfEmpty(9) yielded %v, want [1]", got)
	}
}

var benchmarkDefaultIfEmptySum int

func BenchmarkTypedDefaultIfEmpty(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			legacyFromSlice(source).DefaultIfEmpty(-1).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkDefaultIfEmptySum = sum
		}
	})

	b.Run("generic_query", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).DefaultIfEmpty(-1).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkDefaultIfEmptySum = sum
		}
	})
}
