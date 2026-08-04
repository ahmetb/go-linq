package linq

import (
	"slices"
	"testing"
)

func TestTypedResultsReturnsTypedSliceInOrder(t *testing.T) {
	got := fromSlice([]string{"a", "b", "c"}).Results()
	want := []string{"a", "b", "c"}

	if !slices.Equal(got, want) {
		t.Fatalf("Results() = %v, want %v", got, want)
	}
}

func TestTypedResultsDoesNotAliasSource(t *testing.T) {
	source := []int{1, 2, 3}
	got := fromSlice(source).Results()
	got[0] = 99

	if source[0] != 1 {
		t.Fatalf("Results() mutation changed source to %v", source)
	}
}

func TestTypedResultsReturnsNilForEmptyQuery(t *testing.T) {
	if got := fromSlice([]int{}).Results(); got != nil {
		t.Fatalf("Results() = %v, want nil", got)
	}
}

var (
	benchmarkLegacyResults []any
	benchmarkTypedResults  []int
)

func BenchmarkTypedResults(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkLegacyResults = legacyFromSlice(source).Results()
		}
	})

	b.Run("generic_query", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkTypedResults = fromSlice(source).Results()
		}
	})
}
