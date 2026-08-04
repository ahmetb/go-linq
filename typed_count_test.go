package linq

import "testing"

func TestTypedCountReturnsNumberOfElements(t *testing.T) {
	if got := fromSlice([]string{"a", "b", "c"}).Count(); got != 3 {
		t.Fatalf("Count() = %d, want 3", got)
	}
}

func TestTypedCountReturnsZeroForEmptyQuery(t *testing.T) {
	if got := fromSlice([]int{}).Count(); got != 0 {
		t.Fatalf("Count() = %d, want 0", got)
	}
}

func TestTypedCountEnumeratesEveryElement(t *testing.T) {
	visited := 0
	q := Query[int]{iterate: func(yield func(int) bool) {
		for _, value := range []int{1, 2, 3} {
			visited++
			if !yield(value) {
				return
			}
		}
	}}

	if got := q.Count(); got != 3 {
		t.Fatalf("Count() = %d, want 3", got)
	}
	if visited != 3 {
		t.Fatalf("source visited %d elements, want 3", visited)
	}
}

var benchmarkCountResult int

func BenchmarkTypedCount(b *testing.B) {
	source := make([]int, 1024)

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkCountResult = FromSlice(source).Count()
		}
	})

	b.Run("generic_query", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkCountResult = fromSlice(source).Count()
		}
	})
}
