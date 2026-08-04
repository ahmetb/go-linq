package linq

import (
	"fmt"
	"slices"
	"testing"
)

func TestTypedZipDifferentTypes(t *testing.T) {
	got := fromSlice([]string{"a", "b", "c"}).Zip(
		fromSlice([]int{1, 2}),
		func(left string, right int) string { return fmt.Sprintf("%s%d", left, right) },
	).toSlice()
	want := []string{"a1", "b2"}

	if !slices.Equal(got, want) {
		t.Fatalf("Zip() = %v, want %v", got, want)
	}
}

func TestTypedZipStopsBothInputs(t *testing.T) {
	selectorCalls := 0
	q := fromSlice([]int{1, 2, 3}).Zip(fromSlice([]int{4, 5, 6}), func(left, right int) int {
		selectorCalls++
		return left + right
	})

	q.iterate(func(int) bool { return false })

	if selectorCalls != 1 {
		t.Fatalf("result selector called %d times, want 1", selectorCalls)
	}
}

func zipSum(left, right int) int { return left + right }

var benchmarkZipSum int

func BenchmarkTypedZip(b *testing.B) {
	left := make([]int, 1024)
	right := make([]int, 1024)
	for i := range left {
		left[i], right[i] = i, len(right)-i
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			FromSlice(left).ZipT(FromSlice(right), zipSum).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkZipSum = sum
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(left).Zip(fromSlice(right), zipSum).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkZipSum = sum
		}
	})
}
