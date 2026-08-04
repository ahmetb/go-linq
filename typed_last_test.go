package linq

import "testing"

func TestTypedLastReturnsLastElement(t *testing.T) {
	if got := fromSlice([]string{"first", "last"}).Last(); got != "last" {
		t.Fatalf("Last() = %q, want last", got)
	}
}

func TestTypedLastReturnsZeroForEmptyQuery(t *testing.T) {
	if got := fromSlice([]int{}).Last(); got != 0 {
		t.Fatalf("Last() = %d, want zero value", got)
	}
}

func TestTypedLastEnumeratesEveryElement(t *testing.T) {
	visited := 0
	q := query[int]{iterate: func(yield func(int) bool) {
		for _, value := range []int{1, 2, 3} {
			visited++
			if !yield(value) {
				return
			}
		}
	}}

	if got := q.Last(); got != 3 {
		t.Fatalf("Last() = %d, want 3", got)
	}
	if visited != 3 {
		t.Fatalf("source visited %d elements, want 3", visited)
	}
}

var benchmarkLastResult int

func BenchmarkTypedLast(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkLastResult = FromSlice(source).Last().(int)
		}
	})

	b.Run("generic_query", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkLastResult = fromSlice(source).Last()
		}
	})
}
