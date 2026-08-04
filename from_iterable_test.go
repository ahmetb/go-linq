package linq

import (
	"iter"
	"slices"
	"testing"
)

type typedNumberIterable []int

func (values typedNumberIterable) Iterate() iter.Seq[int] {
	return slices.Values(values)
}

type legacyNumberIterable []int

func (values legacyNumberIterable) Iterate() iter.Seq[any] {
	return func(yield func(any) bool) {
		for _, value := range values {
			if !yield(value) {
				return
			}
		}
	}
}

func TestTypedFromIterablePreservesElementType(t *testing.T) {
	got := fromIterable(typedNumberIterable{1, 2, 3}).toSlice()
	if !slices.Equal(got, []int{1, 2, 3}) {
		t.Fatalf("fromIterable() = %v, want [1 2 3]", got)
	}
}

func TestTypedFromIterableStopsWithConsumer(t *testing.T) {
	var got []int
	fromIterable(typedNumberIterable{1, 2, 3}).iterate(func(value int) bool {
		got = append(got, value)
		return false
	})

	if !slices.Equal(got, []int{1}) {
		t.Fatalf("fromIterable() yielded %v, want [1]", got)
	}
}

var benchmarkFromIterableSum int

func BenchmarkTypedFromIterable(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			legacyFromIterable(legacyNumberIterable(source)).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkFromIterableSum = sum
		}
	})

	b.Run("generic_query", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromIterable(typedNumberIterable(source)).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkFromIterableSum = sum
		}
	})
}
