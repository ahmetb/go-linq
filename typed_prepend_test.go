package linq

import (
	"slices"
	"testing"
)

func TestTypedPrependAddsValueBeforeSource(t *testing.T) {
	got := fromSlice([]int{2, 3}).Prepend(1).toSlice()
	if !slices.Equal(got, []int{1, 2, 3}) {
		t.Fatalf("Prepend(1) = %v, want [1 2 3]", got)
	}
}

func TestTypedPrependAddsValueToEmptyQuery(t *testing.T) {
	got := fromSlice([]string{}).Prepend("only").toSlice()
	if !slices.Equal(got, []string{"only"}) {
		t.Fatalf("Prepend(only) = %v, want [only]", got)
	}
}

func TestTypedPrependDoesNotEnumerateSourceAfterConsumerStops(t *testing.T) {
	visited := 0
	source := Query[int]{iterate: func(yield func(int) bool) {
		visited++
		yield(2)
	}}

	source.Prepend(1).iterate(func(int) bool { return false })
	if visited != 0 {
		t.Fatalf("source visited %d elements, want 0", visited)
	}
}

var benchmarkPrependSum int

func BenchmarkTypedPrepend(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			FromSlice(source).Prepend(-1).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkPrependSum = sum
		}
	})

	b.Run("generic_query", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).Prepend(-1).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkPrependSum = sum
		}
	})
}
