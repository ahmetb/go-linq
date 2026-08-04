package linq

import "testing"

func TestTypedSumByAcceptsMethodExpression(t *testing.T) {
	people := []joinPerson{{id: 2}, {id: 5}, {id: 3}}
	got := fromSlice(people).SumBy(joinPerson.ID)

	if got != 10 {
		t.Fatalf("SumBy(joinPerson.ID) = %d, want 10", got)
	}
}

func TestTypedSumByPreservesNamedNumericType(t *testing.T) {
	type score int16
	got := fromSlice([]score{1, 2, 3}).SumBy(func(value score) score { return value })

	var want score = 6
	if got != want {
		t.Fatalf("SumBy() = %d, want %d", got, want)
	}
}

func TestTypedSumByReturnsZeroForEmptyQuery(t *testing.T) {
	got := fromSlice([]string{}).SumBy(func(string) float64 { return 1 })
	if got != 0 {
		t.Fatalf("SumBy() = %v, want 0", got)
	}
}

func sumByBenchmarkValue(value int) int { return value }

var benchmarkSumByResult int64

func BenchmarkTypedSumBy(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkSumByResult = FromSlice(source).SumInts()
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkSumByResult = int64(fromSlice(source).SumBy(sumByBenchmarkValue))
		}
	})
}
