package linq

import (
	"slices"
	"testing"
)

func TestTypedIntersectByMethodExpressionReturnsSetIntersection(t *testing.T) {
	first := joinPerson{id: 1, name: "first"}
	source := []joinPerson{first, {id: 1, name: "duplicate"}, {id: 2, name: "excluded"}, {id: 3, name: "third"}}
	other := []joinPerson{{id: 1}, {id: 3}}

	got := fromSlice(source).IntersectBy(fromSlice(other), joinPerson.ID).toSlice()
	want := []joinPerson{first, source[3]}

	if !slices.Equal(got, want) {
		t.Fatalf("IntersectBy(joinPerson.ID) = %v, want %v", got, want)
	}
}

func TestTypedIntersectByStopsWithConsumer(t *testing.T) {
	calls := 0
	q := fromSlice([]int{1, 2, 3}).IntersectBy(fromSlice([]int{1, 2}), func(value int) int {
		calls++
		return value
	})

	q.iterate(func(int) bool { return false })

	if calls != 3 {
		t.Fatalf("selector called %d times, want 3 (right side plus first matching left value)", calls)
	}
}

func intersectBenchmarkKey(value int) int { return value }

var benchmarkIntersectBySum int

func BenchmarkTypedIntersectBy(b *testing.B) {
	source := make([]int, 1024)
	other := make([]int, 512)
	for i := range source {
		source[i] = i
	}
	copy(other, source[:len(other)])

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			FromSlice(source).IntersectByT(FromSlice(other), intersectBenchmarkKey).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkIntersectBySum = sum
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).IntersectBy(fromSlice(other), intersectBenchmarkKey).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkIntersectBySum = sum
		}
	})
}
