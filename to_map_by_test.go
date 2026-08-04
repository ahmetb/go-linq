package linq

import (
	"maps"
	"testing"
)

func joinPersonName(person joinPerson) string { return person.name }

func TestTypedToMapByMethodExpressionsReturnTypedMap(t *testing.T) {
	people := []joinPerson{
		{id: 1, name: "first"},
		{id: 2, name: "second"},
		{id: 1, name: "replacement"},
	}

	got := fromSlice(people).ToMapBy(joinPerson.ID, joinPersonName)
	want := map[int]string{1: "replacement", 2: "second"}

	if !maps.Equal(got, want) {
		t.Fatalf("ToMapBy(joinPerson.ID, joinPersonName) = %v, want %v", got, want)
	}
}

func TestTypedToMapByReturnsNonNilEmptyMap(t *testing.T) {
	got := fromSlice([]int{}).ToMapBy(func(value int) int { return value }, func(value int) int { return value })

	if got == nil {
		t.Fatal("ToMapBy() returned a nil map for an empty query")
	}
}

func toMapBenchmarkKey(value int) int { return value % 128 }

func toMapBenchmarkValue(value int) int { return value * 2 }

var benchmarkToMapByResult map[int]int

func BenchmarkTypedToMapBy(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			result := make(map[int]int)
			legacyFromSlice(source).ToMapByT(&result, toMapBenchmarkKey, toMapBenchmarkValue)
			benchmarkToMapByResult = result
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkToMapByResult = fromSlice(source).ToMapBy(toMapBenchmarkKey, toMapBenchmarkValue)
		}
	})
}
