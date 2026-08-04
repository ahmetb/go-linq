package linq

import (
	"slices"
	"testing"
)

type selectManyWidget struct {
	values []int
}

func (w *selectManyWidget) Values() Query[int] {
	return fromSlice(w.values)
}

func TestTypedSelectManyMethodExpression(t *testing.T) {
	source := []*selectManyWidget{{values: []int{1, 2}}, {values: nil}, {values: []int{3, 4}}}

	got := fromSlice(source).
		SelectMany((*selectManyWidget).Values).
		toSlice()
	want := []int{1, 2, 3, 4}

	if !slices.Equal(got, want) {
		t.Fatalf("SelectMany((*selectManyWidget).Values) = %v, want %v", got, want)
	}
}

func TestTypedSelectManyStopsBothIterators(t *testing.T) {
	selectorCalls := 0
	q := fromSlice([][]int{{1, 2}, {3, 4}}).SelectMany(func(values []int) Query[int] {
		selectorCalls++
		return fromSlice(values)
	})

	var got []int
	q.iterate(func(value int) bool {
		got = append(got, value)
		return false
	})

	if want := []int{1}; !slices.Equal(got, want) {
		t.Fatalf("SelectMany iteration = %v, want %v", got, want)
	}
	if selectorCalls != 1 {
		t.Fatalf("selector called %d times, want 1", selectorCalls)
	}
}

func selectManyLegacy(values []int) legacyQuery { return legacyFromSlice(values) }

func selectManyTyped(values []int) Query[int] { return fromSlice(values) }

var benchmarkSelectManySum int

func BenchmarkTypedSelectMany(b *testing.B) {
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
			legacyFromSlice(source).SelectManyT(selectManyLegacy).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkSelectManySum = sum
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).SelectMany(selectManyTyped).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkSelectManySum = sum
		}
	})
}
