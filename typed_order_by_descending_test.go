package linq

import (
	"slices"
	"testing"
)

func TestTypedOrderByDescendingMethodExpressionIsStable(t *testing.T) {
	source := []orderWidget{{rank: 1, name: "a"}, {rank: 2, name: "c"}, {rank: 2, name: "d"}}

	got := fromSlice(source).OrderByDescending(orderWidget.Rank).toSlice()
	want := []orderWidget{{rank: 2, name: "c"}, {rank: 2, name: "d"}, {rank: 1, name: "a"}}

	if !slices.Equal(got, want) {
		t.Fatalf("OrderByDescending(orderWidget.Rank) = %v, want %v", got, want)
	}
}

var benchmarkOrderByDescendingSum int

func BenchmarkTypedOrderByDescending(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = (i * 353) % len(source)
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			legacyFromSlice(source).OrderByDescendingT(orderIdentity).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkOrderByDescendingSum = sum
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).OrderByDescending(orderIdentity).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkOrderByDescendingSum = sum
		}
	})
}
