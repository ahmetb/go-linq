package linq

import (
	"slices"
	"testing"
)

type forEachValue struct {
	visited bool
}

func (value *forEachValue) Visit() { value.visited = true }

func TestTypedForEachAcceptsMethodExpression(t *testing.T) {
	values := []*forEachValue{{}, {}}
	fromSlice(values).ForEach((*forEachValue).Visit)

	for i, value := range values {
		if !value.visited {
			t.Fatalf("value %d was not visited", i)
		}
	}
}

func TestTypedForEachPreservesOrder(t *testing.T) {
	var got []int
	fromSlice([]int{3, 1, 2}).ForEach(func(value int) {
		got = append(got, value)
	})

	want := []int{3, 1, 2}
	if !slices.Equal(got, want) {
		t.Fatalf("ForEach() visited %v, want %v", got, want)
	}
}

var benchmarkForEachSum int

func BenchmarkTypedForEach(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			FromSlice(source).ForEachT(func(value int) {
				sum += value
			})
			benchmarkForEachSum = sum
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).ForEach(func(value int) {
				sum += value
			})
			benchmarkForEachSum = sum
		}
	})
}
