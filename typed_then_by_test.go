package linq

import (
	"slices"
	"testing"
)

func (w orderWidget) Name() string { return w.name }

func TestTypedThenBySupportsDifferentKeyType(t *testing.T) {
	source := []orderWidget{
		{rank: 2, name: "a"}, {rank: 1, name: "c"},
		{rank: 1, name: "b"}, {rank: 2, name: "b"},
	}

	got := fromSlice(source).
		OrderBy(orderWidget.Rank).
		ThenBy(orderWidget.Name).
		toSlice()
	want := []orderWidget{
		{rank: 1, name: "b"}, {rank: 1, name: "c"},
		{rank: 2, name: "a"}, {rank: 2, name: "b"},
	}

	if !slices.Equal(got, want) {
		t.Fatalf("OrderBy().ThenBy() = %v, want %v", got, want)
	}
}

func TestTypedThenByDoesNotMutateParentOrdering(t *testing.T) {
	source := []orderWidget{{rank: 1, name: "b"}, {rank: 1, name: "a"}}
	primary := fromSlice(source).OrderBy(orderWidget.Rank)
	_ = primary.ThenBy(orderWidget.Name)

	if got := primary.toSlice(); !slices.Equal(got, source) {
		t.Fatalf("parent ordering mutated: got %v, want %v", got, source)
	}
}

var benchmarkThenBySum int

func BenchmarkTypedThenBy(b *testing.B) {
	source := make([]orderWidget, 1024)
	for i := range source {
		source[i] = orderWidget{rank: i % 32, name: string(rune('a' + (i*17)%26))}
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			legacyFromSlice(source).OrderByT(orderWidget.Rank).ThenByT(orderWidget.Name).Iterate(func(value any) bool {
				sum += value.(orderWidget).rank
				return true
			})
			benchmarkThenBySum = sum
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).OrderBy(orderWidget.Rank).ThenBy(orderWidget.Name).iterate(func(value orderWidget) bool {
				sum += value.rank
				return true
			})
			benchmarkThenBySum = sum
		}
	})
}
