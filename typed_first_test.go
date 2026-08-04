package linq

import "testing"

func TestTypedFirstReturnsFirstElement(t *testing.T) {
	if got := fromSlice([]string{"first", "second"}).First(); got != "first" {
		t.Fatalf("First() = %q, want first", got)
	}
}

func TestTypedFirstReturnsZeroForEmptyQuery(t *testing.T) {
	if got := fromSlice([]int{}).First(); got != 0 {
		t.Fatalf("First() = %d, want zero value", got)
	}
}

func TestTypedFirstStopsAfterFirstElement(t *testing.T) {
	visited := 0
	q := query[int]{iterate: func(yield func(int) bool) {
		for _, value := range []int{1, 2, 3} {
			visited++
			if !yield(value) {
				return
			}
		}
	}}

	if got := q.First(); got != 1 {
		t.Fatalf("First() = %d, want 1", got)
	}
	if visited != 1 {
		t.Fatalf("source visited %d elements, want 1", visited)
	}
}

var benchmarkFirstResult int

func BenchmarkTypedFirst(b *testing.B) {
	source := make([]int, 1024)
	source[0] = 42

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkFirstResult = FromSlice(source).First().(int)
		}
	})

	b.Run("generic_query", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkFirstResult = fromSlice(source).First()
		}
	})
}
