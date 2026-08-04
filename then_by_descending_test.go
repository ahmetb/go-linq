package linq

import (
	"slices"
	"testing"
)

func TestTypedThenByDescendingSupportsDifferentKeyType(t *testing.T) {
	source := []orderWidget{
		{rank: 2, name: "a"}, {rank: 1, name: "b"},
		{rank: 1, name: "c"}, {rank: 2, name: "b"},
	}

	got := fromSlice(source).
		OrderBy(orderWidget.Rank).
		ThenByDescending(orderWidget.Name).
		toSlice()
	want := []orderWidget{
		{rank: 1, name: "c"}, {rank: 1, name: "b"},
		{rank: 2, name: "b"}, {rank: 2, name: "a"},
	}

	if !slices.Equal(got, want) {
		t.Fatalf("OrderBy().ThenByDescending() = %v, want %v", got, want)
	}
}

var benchmarkThenByDescendingSum int

func BenchmarkTypedThenByDescending(b *testing.B) {
	source := make([]orderWidget, 1024)
	for i := range source {
		source[i] = orderWidget{rank: i % 32, name: string(rune('a' + (i*17)%26))}
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			legacyFromSlice(source).OrderByT(orderWidget.Rank).ThenByDescendingT(orderWidget.Name).Iterate(func(value any) bool {
				sum += value.(orderWidget).rank
				return true
			})
			benchmarkThenByDescendingSum = sum
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).OrderBy(orderWidget.Rank).ThenByDescending(orderWidget.Name).iterate(func(value orderWidget) bool {
				sum += value.rank
				return true
			})
			benchmarkThenByDescendingSum = sum
		}
	})
}
