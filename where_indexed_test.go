package linq

import (
	"slices"
	"testing"
)

func TestTypedWhereIndexed(t *testing.T) {
	got := fromSlice([]string{"zero", "one", "two", "three"}).
		WhereIndexed(func(index int, _ string) bool { return index%2 == 0 }).
		toSlice()
	want := []string{"zero", "two"}

	if !slices.Equal(got, want) {
		t.Fatalf("WhereIndexed() = %v, want %v", got, want)
	}
}

func TestTypedWhereIndexedStopsAfterFirstAcceptedValue(t *testing.T) {
	var indices []int
	q := fromSlice([]int{-1, -1, 1, 1}).WhereIndexed(func(index, value int) bool {
		indices = append(indices, index)
		return value > 0
	})

	q.iterate(func(int) bool { return false })

	if want := []int{0, 1, 2}; !slices.Equal(indices, want) {
		t.Fatalf("predicate indices = %v, want %v", indices, want)
	}
}

func whereIndexedEven(index, _ int) bool { return index%2 == 0 }

var benchmarkWhereIndexedSum int

func BenchmarkTypedWhereIndexed(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			legacyFromSlice(source).WhereIndexedT(whereIndexedEven).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkWhereIndexedSum = sum
		}
	})

	b.Run("generic_query", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).WhereIndexed(whereIndexedEven).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkWhereIndexedSum = sum
		}
	})
}
