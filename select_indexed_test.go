package linq

import (
	"fmt"
	"slices"
	"testing"
)

func TestTypedSelectIndexed(t *testing.T) {
	got := fromSlice([]string{"a", "b", "c"}).
		SelectIndexed(func(index int, value string) string {
			return fmt.Sprintf("%d:%s", index, value)
		}).
		toSlice()
	want := []string{"0:a", "1:b", "2:c"}

	if !slices.Equal(got, want) {
		t.Fatalf("SelectIndexed() = %v, want %v", got, want)
	}
}

func TestTypedSelectIndexedStopsWithConsumer(t *testing.T) {
	var indices []int
	q := fromSlice([]string{"a", "b", "c"}).SelectIndexed(func(index int, value string) string {
		indices = append(indices, index)
		return value
	})

	q.iterate(func(string) bool { return false })

	if want := []int{0}; !slices.Equal(indices, want) {
		t.Fatalf("selector indices = %v, want %v", indices, want)
	}
}

func selectIndexedSum(index, value int) int { return index + value }

var benchmarkSelectIndexedSum int

func BenchmarkTypedSelectIndexed(b *testing.B) {
	source := make([]int, 1024)

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			legacyFromSlice(source).SelectIndexedT(selectIndexedSum).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkSelectIndexedSum = sum
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).SelectIndexed(selectIndexedSum).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkSelectIndexedSum = sum
		}
	})
}
