package linq

import (
	"slices"
	"testing"
)

func TestTypedForEachIndexedSuppliesIndexAndValue(t *testing.T) {
	var got []int
	fromSlice([]int{10, 20, 30}).ForEachIndexed(func(index, value int) {
		got = append(got, index+value)
	})

	want := []int{10, 21, 32}
	if !slices.Equal(got, want) {
		t.Fatalf("ForEachIndexed() produced %v, want %v", got, want)
	}
}

func TestTypedForEachIndexedDoesNotRunForEmptyQuery(t *testing.T) {
	calls := 0
	fromSlice([]string{}).ForEachIndexed(func(int, string) { calls++ })

	if calls != 0 {
		t.Fatalf("action called %d times, want 0", calls)
	}
}

var benchmarkForEachIndexedSum int

func BenchmarkTypedForEachIndexed(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			legacyFromSlice(source).ForEachIndexedT(func(index, value int) {
				sum += index + value
			})
			benchmarkForEachIndexedSum = sum
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).ForEachIndexed(func(index, value int) {
				sum += index + value
			})
			benchmarkForEachIndexedSum = sum
		}
	})
}
