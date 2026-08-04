package linq

import (
	"math"
	"testing"
)

func TestTypedAverageByAcceptsMethodExpression(t *testing.T) {
	people := []joinPerson{{id: 2}, {id: 5}, {id: 8}}
	got := fromSlice(people).AverageBy(joinPerson.ID)

	if got != 5 {
		t.Fatalf("AverageBy(joinPerson.ID) = %v, want 5", got)
	}
}

func TestTypedAverageBySupportsNamedFloats(t *testing.T) {
	type measurement float32
	got := fromSlice([]measurement{1, 2, 3}).AverageBy(func(value measurement) measurement { return value })

	if got != 2 {
		t.Fatalf("AverageBy() = %v, want 2", got)
	}
}

func TestTypedAverageByReturnsNaNForEmptyQuery(t *testing.T) {
	got := fromSlice([]int{}).AverageBy(func(value int) int { return value })
	if !math.IsNaN(got) {
		t.Fatalf("AverageBy() = %v, want NaN", got)
	}
}

func averageByBenchmarkValue(value int) int { return value }

var benchmarkAverageByResult float64

func BenchmarkTypedAverageBy(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkAverageByResult = FromSlice(source).Average()
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkAverageByResult = fromSlice(source).AverageBy(averageByBenchmarkValue)
		}
	})
}
