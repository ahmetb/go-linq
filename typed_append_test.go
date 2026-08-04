package linq

import (
	"slices"
	"testing"
)

func TestTypedAppendAddsValueAfterSource(t *testing.T) {
	got := fromSlice([]int{1, 2}).Append(3).toSlice()
	if !slices.Equal(got, []int{1, 2, 3}) {
		t.Fatalf("Append(3) = %v, want [1 2 3]", got)
	}
}

func TestTypedAppendAddsValueToEmptyQuery(t *testing.T) {
	got := fromSlice([]string{}).Append("only").toSlice()
	if !slices.Equal(got, []string{"only"}) {
		t.Fatalf("Append(only) = %v, want [only]", got)
	}
}

func TestTypedAppendDoesNotReachValueAfterConsumerStops(t *testing.T) {
	var got []int
	fromSlice([]int{1, 2}).Append(3).iterate(func(value int) bool {
		got = append(got, value)
		return false
	})

	if !slices.Equal(got, []int{1}) {
		t.Fatalf("Append() yielded %v, want [1]", got)
	}
}

var benchmarkAppendSum int

func BenchmarkTypedAppend(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			FromSlice(source).Append(1024).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkAppendSum = sum
		}
	})

	b.Run("generic_query", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).Append(1024).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkAppendSum = sum
		}
	})
}
