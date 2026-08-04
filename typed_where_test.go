package linq

import (
	"slices"
	"testing"
)

type whereWidget struct {
	enabled bool
}

func (w *whereWidget) Enabled() bool {
	return w.enabled
}

func whereEven(value int) bool {
	return value%2 == 0
}

func whereEvenAny(value any) bool {
	return value.(int)%2 == 0
}

func TestTypedWhereMethodExpression(t *testing.T) {
	source := []*whereWidget{{enabled: true}, {enabled: false}, {enabled: true}}

	got := fromSlice(source).
		Where((*whereWidget).Enabled).
		toSlice()
	want := []*whereWidget{source[0], source[2]}

	if !slices.Equal(got, want) {
		t.Fatalf("Where((*whereWidget).Enabled) = %v, want %v", got, want)
	}
}

func TestTypedWhereStopsAfterFirstAcceptedValue(t *testing.T) {
	calls := 0
	q := fromSlice([]int{-1, 0, 2, 3}).Where(func(value int) bool {
		calls++
		return value > 0
	})

	var got []int
	q.iterate(func(value int) bool {
		got = append(got, value)
		return false
	})

	if want := []int{2}; !slices.Equal(got, want) {
		t.Fatalf("Where iteration = %v, want %v", got, want)
	}
	if calls != 3 {
		t.Fatalf("predicate called %d times, want 3", calls)
	}
}

var benchmarkWhereSum int

func BenchmarkTypedWhere(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy_any", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			FromSlice(source).Where(whereEvenAny).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkWhereSum = sum
		}
	})

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			FromSlice(source).WhereT(whereEven).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkWhereSum = sum
		}
	})

	b.Run("generic_query", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).Where(whereEven).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkWhereSum = sum
		}
	})
}
