package linq

import (
	"slices"
	"testing"
)

type selectWidget struct {
	value int
}

func (w *selectWidget) Value() int {
	return w.value
}

func selectDouble(value int) int {
	return value * 2
}

func selectDoubleAny(value any) any {
	return value.(int) * 2
}

func TestTypedSelectMethodExpression(t *testing.T) {
	source := []*selectWidget{{value: 1}, {value: 2}, {value: 3}}

	got := fromSlice(source).
		Select((*selectWidget).Value).
		toSlice()
	want := []int{1, 2, 3}

	if !slices.Equal(got, want) {
		t.Fatalf("Select((*selectWidget).Value) = %v, want %v", got, want)
	}
}

func TestTypedSelectStopsWhenConsumerStops(t *testing.T) {
	calls := 0
	q := fromSlice([]int{1, 2, 3}).Select(func(value int) int {
		calls++
		return value * 2
	})

	q.iterate(func(int) bool { return false })

	if calls != 1 {
		t.Fatalf("selector called %d times, want 1", calls)
	}
}

var benchmarkSelectSum int

func BenchmarkTypedSelect(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy_any", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			FromSlice(source).Select(selectDoubleAny).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkSelectSum = sum
		}
	})

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			FromSlice(source).SelectT(selectDouble).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkSelectSum = sum
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).Select(selectDouble).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkSelectSum = sum
		}
	})
}
