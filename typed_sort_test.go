package linq

import (
	"slices"
	"testing"
)

type sortNumber int

func (number sortNumber) Less(other sortNumber) bool { return number < other }

func TestTypedSortAcceptsMethodExpression(t *testing.T) {
	source := []sortNumber{3, 1, 2}
	got := fromSlice(source).Sort(sortNumber.Less).toSlice()
	want := []sortNumber{1, 2, 3}

	if !slices.Equal(got, want) {
		t.Fatalf("Sort(sortNumber.Less) = %v, want %v", got, want)
	}
	if !slices.Equal(source, []sortNumber{3, 1, 2}) {
		t.Fatalf("Sort() mutated source to %v", source)
	}
}

func TestTypedSortDefersWorkUntilEnumeration(t *testing.T) {
	calls := 0
	q := fromSlice([]int{2, 1}).Sort(func(left, right int) bool {
		calls++
		return left < right
	})

	if calls != 0 {
		t.Fatalf("comparator called %d times before enumeration", calls)
	}
	q.toSlice()
	if calls == 0 {
		t.Fatal("comparator was not called during enumeration")
	}
}

func sortBenchmarkLess(left, right int) bool { return left < right }

var benchmarkSortSum int

func BenchmarkTypedSort(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = len(source) - i
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			legacyFromSlice(source).SortT(sortBenchmarkLess).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkSortSum = sum
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).Sort(sortBenchmarkLess).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkSortSum = sum
		}
	})
}
