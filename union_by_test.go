package linq

import (
	"slices"
	"testing"
)

func TestTypedUnionByAcceptsMethodExpressionAndPreservesFirstValues(t *testing.T) {
	left := []joinPerson{{id: 1, name: "left one"}, {id: 1, name: "duplicate"}, {id: 2, name: "left two"}}
	right := []joinPerson{{id: 2, name: "right duplicate"}, {id: 3, name: "right three"}}

	got := fromSlice(left).UnionBy(fromSlice(right), joinPerson.ID).toSlice()
	want := []joinPerson{left[0], left[2], right[1]}
	if !slices.Equal(got, want) {
		t.Fatalf("UnionBy(joinPerson.ID) = %v, want %v", got, want)
	}
}

func TestTypedUnionByDoesNotEnumerateRightAfterConsumerStops(t *testing.T) {
	rightVisited := 0
	right := Query[int]{iterate: func(yield func(int) bool) {
		rightVisited++
		yield(2)
	}}

	fromSlice([]int{1}).UnionBy(right, func(value int) int { return value }).iterate(func(int) bool { return false })
	if rightVisited != 0 {
		t.Fatalf("right query visited %d elements, want 0", rightVisited)
	}
}

func unionByBenchmarkKey(value int) int { return value }

var benchmarkUnionBySum int

func BenchmarkTypedUnionBy(b *testing.B) {
	left := make([]int, 1024)
	right := make([]int, 1024)
	for i := range left {
		left[i] = i % 512
		right[i] = 256 + i%512
	}

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			legacyFromSlice(left).Union(legacyFromSlice(right)).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkUnionBySum = sum
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(left).UnionBy(fromSlice(right), unionByBenchmarkKey).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkUnionBySum = sum
		}
	})
}
