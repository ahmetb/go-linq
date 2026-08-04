package linq

import (
	"slices"
	"testing"
)

func TestTypedSelectManyIndexed(t *testing.T) {
	got := fromSlice([]string{"a", "b", "c"}).
		SelectManyIndexed(func(index int, value string) query[string] {
			return fromSlice([]string{value, value + string(rune('0'+index))})
		}).
		toSlice()
	want := []string{"a", "a0", "b", "b1", "c", "c2"}

	if !slices.Equal(got, want) {
		t.Fatalf("SelectManyIndexed() = %v, want %v", got, want)
	}
}

func TestTypedSelectManyIndexedStopsBothIterators(t *testing.T) {
	var indices []int
	q := fromSlice([][]int{{1, 2}, {3, 4}}).SelectManyIndexed(func(index int, values []int) query[int] {
		indices = append(indices, index)
		return fromSlice(values)
	})

	q.iterate(func(int) bool { return false })

	if want := []int{0}; !slices.Equal(indices, want) {
		t.Fatalf("selector indices = %v, want %v", indices, want)
	}
}

func selectManyIndexedLegacy(_ int, values []int) Query { return FromSlice(values) }

func selectManyIndexedTyped(_ int, values []int) query[int] { return fromSlice(values) }

var benchmarkSelectManyIndexedSum int

func BenchmarkTypedSelectManyIndexed(b *testing.B) {
	source := make([][]int, 128)
	for i := range source {
		source[i] = make([]int, 8)
		for j := range source[i] {
			source[i][j] = i*len(source[i]) + j
		}
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			FromSlice(source).SelectManyIndexedT(selectManyIndexedLegacy).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkSelectManyIndexedSum = sum
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).SelectManyIndexed(selectManyIndexedTyped).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkSelectManyIndexedSum = sum
		}
	})
}
