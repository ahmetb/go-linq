package linq

import (
	"slices"
	"testing"
)

func TestTypedConcatPreservesSequenceOrder(t *testing.T) {
	got := fromSlice([]int{1, 2}).Concat(fromSlice([]int{3, 4})).toSlice()
	if !slices.Equal(got, []int{1, 2, 3, 4}) {
		t.Fatalf("Concat() = %v, want [1 2 3 4]", got)
	}
}

func TestTypedConcatHandlesEmptyQueries(t *testing.T) {
	got := fromSlice([]string{}).Concat(fromSlice([]string{"value"})).toSlice()
	if !slices.Equal(got, []string{"value"}) {
		t.Fatalf("Concat() = %v, want [value]", got)
	}
}

func TestTypedConcatDoesNotEnumerateRightAfterConsumerStops(t *testing.T) {
	rightVisited := 0
	right := query[int]{iterate: func(yield func(int) bool) {
		rightVisited++
		yield(2)
	}}

	fromSlice([]int{1}).Concat(right).iterate(func(int) bool { return false })
	if rightVisited != 0 {
		t.Fatalf("right query visited %d elements, want 0", rightVisited)
	}
}

var benchmarkConcatSum int

func BenchmarkTypedConcat(b *testing.B) {
	left := make([]int, 512)
	right := make([]int, 512)
	for i := range left {
		left[i] = i
		right[i] = i + len(left)
	}

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			FromSlice(left).Concat(FromSlice(right)).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkConcatSum = sum
		}
	})

	b.Run("generic_query", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(left).Concat(fromSlice(right)).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkConcatSum = sum
		}
	})
}
