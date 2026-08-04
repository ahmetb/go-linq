package linq

import (
	"slices"
	"testing"
)

func TestTypedReverseReturnsValuesInReverseOrder(t *testing.T) {
	source := []int{1, 2, 3}
	got := fromSlice(source).Reverse().toSlice()

	if !slices.Equal(got, []int{3, 2, 1}) {
		t.Fatalf("Reverse() = %v, want [3 2 1]", got)
	}
	if !slices.Equal(source, []int{1, 2, 3}) {
		t.Fatalf("Reverse() mutated source to %v", source)
	}
}

func TestTypedReverseStopsWithConsumer(t *testing.T) {
	var got []string
	fromSlice([]string{"a", "b", "c"}).Reverse().iterate(func(value string) bool {
		got = append(got, value)
		return false
	})

	if !slices.Equal(got, []string{"c"}) {
		t.Fatalf("Reverse() yielded %v, want [c]", got)
	}
}

var benchmarkReverseSum int

func BenchmarkTypedReverse(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			FromSlice(source).Reverse().Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkReverseSum = sum
		}
	})

	b.Run("generic_query", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).Reverse().iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkReverseSum = sum
		}
	})
}
