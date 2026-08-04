package linq

import (
	"slices"
	"testing"
)

func TestTypedSelectManyByIndexed(t *testing.T) {
	got := fromSlice([]string{"a", "b"}).SelectManyByIndexed(
		func(index int, value string) Query[int] {
			return fromSlice([]int{index, len(value)})
		},
		func(inner int, outer string) string { return outer + string(rune('0'+inner)) },
	).toSlice()
	want := []string{"a0", "a1", "b1", "b1"}

	if !slices.Equal(got, want) {
		t.Fatalf("SelectManyByIndexed() = %v, want %v", got, want)
	}
}

func TestTypedSelectManyByIndexedStopsBothIterators(t *testing.T) {
	var indices []int
	q := fromSlice([][]int{{1, 2}, {3, 4}}).SelectManyByIndexed(
		func(index int, values []int) Query[int] {
			indices = append(indices, index)
			return fromSlice(values)
		},
		func(inner int, outer []int) int { return inner + len(outer) },
	)

	q.iterate(func(int) bool { return false })

	if want := []int{0}; !slices.Equal(indices, want) {
		t.Fatalf("selector indices = %v, want %v", indices, want)
	}
}

var benchmarkSelectManyByIndexedSum int

func BenchmarkTypedSelectManyByIndexed(b *testing.B) {
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
			FromSlice(source).SelectManyByIndexedT(selectManyIndexedLegacy, selectManyResult).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkSelectManyByIndexedSum = sum
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).SelectManyByIndexed(selectManyIndexedTyped, selectManyResult).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkSelectManyByIndexedSum = sum
		}
	})
}
