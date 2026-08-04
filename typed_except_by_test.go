package linq

import (
	"slices"
	"testing"
)

func TestTypedExceptByMethodExpressionReturnsSetDifference(t *testing.T) {
	first := joinPerson{id: 1, name: "first"}
	source := []joinPerson{first, {id: 1, name: "duplicate"}, {id: 2, name: "excluded"}, {id: 3, name: "third"}}
	other := []joinPerson{{id: 2}}

	got := fromSlice(source).ExceptBy(fromSlice(other), joinPerson.ID).toSlice()
	want := []joinPerson{first, source[3]}

	if !slices.Equal(got, want) {
		t.Fatalf("ExceptBy(joinPerson.ID) = %v, want %v", got, want)
	}
}

func TestTypedExceptByStopsWithConsumer(t *testing.T) {
	calls := 0
	q := fromSlice([]int{1, 2, 3}).ExceptBy(fromSlice([]int{2}), func(value int) int {
		calls++
		return value
	})

	q.iterate(func(int) bool { return false })

	if calls != 2 {
		t.Fatalf("selector called %d times, want 2 (right side plus first accepted left value)", calls)
	}
}

func exceptBenchmarkKey(value int) int { return value }

var benchmarkExceptBySum int

func BenchmarkTypedExceptBy(b *testing.B) {
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
			legacyFromSlice(source).ExceptByT(legacyFromSlice(other), exceptBenchmarkKey).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkExceptBySum = sum
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).ExceptBy(fromSlice(other), exceptBenchmarkKey).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkExceptBySum = sum
		}
	})
}
