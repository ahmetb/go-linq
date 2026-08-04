package linq

import (
	"maps"
	"testing"
)

func TestTypedFromMapPreservesKeyAndValueTypes(t *testing.T) {
	source := map[int]string{1: "one", 2: "two"}
	got := fromMap(source).ToMapBy(
		func(item KeyValue[int, string]) int { return item.Key },
		func(item KeyValue[int, string]) string { return item.Value },
	)

	if !maps.Equal(got, source) {
		t.Fatalf("fromMap() = %v, want %v", got, source)
	}
}

func TestTypedFromMapStopsWithConsumer(t *testing.T) {
	calls := 0
	fromMap(map[int]int{1: 1, 2: 2}).iterate(func(KeyValue[int, int]) bool {
		calls++
		return false
	})

	if calls != 1 {
		t.Fatalf("consumer called %d times, want 1", calls)
	}
}

var benchmarkFromMapSum int

func BenchmarkTypedFromMap(b *testing.B) {
	source := make(map[int]int, 1024)
	for i := 0; i < 1024; i++ {
		source[i] = i * 2
	}

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			legacyFromMap(source).Iterate(func(value any) bool {
				item := value.(legacyKeyValue)
				sum += item.Key.(int) + item.Value.(int)
				return true
			})
			benchmarkFromMapSum = sum
		}
	})

	b.Run("generic_query", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromMap(source).iterate(func(item KeyValue[int, int]) bool {
				sum += item.Key + item.Value
				return true
			})
			benchmarkFromMapSum = sum
		}
	})
}
