package linq

import (
	"slices"
	"testing"
)

func TestTypedDistinctByMethodExpression(t *testing.T) {
	source := []joinPerson{{id: 1, name: "first"}, {id: 1, name: "duplicate"}, {id: 2, name: "second"}}

	got := fromSlice(source).DistinctBy(joinPerson.ID).toSlice()
	want := []joinPerson{source[0], source[2]}

	if !slices.Equal(got, want) {
		t.Fatalf("DistinctBy(joinPerson.ID) = %v, want %v", got, want)
	}
}

func TestTypedDistinctByStopsWithConsumer(t *testing.T) {
	calls := 0
	q := fromSlice([]int{1, 1, 2}).DistinctBy(func(value int) int {
		calls++
		return value
	})

	q.iterate(func(int) bool { return false })

	if calls != 1 {
		t.Fatalf("selector called %d times, want 1", calls)
	}
}

func distinctBenchmarkKey(value int) int { return value % 128 }

var benchmarkDistinctBySum int

func BenchmarkTypedDistinctBy(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			legacyFromSlice(source).DistinctByT(distinctBenchmarkKey).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkDistinctBySum = sum
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).DistinctBy(distinctBenchmarkKey).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkDistinctBySum = sum
		}
	})
}
