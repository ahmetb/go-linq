package linq

import "testing"

func TestTypedSingleReturnsOnlyElement(t *testing.T) {
	if got := fromSlice([]string{"only"}).Single(); got != "only" {
		t.Fatalf("Single() = %q, want only", got)
	}
}

func TestTypedSingleReturnsZeroUnlessExactlyOneElement(t *testing.T) {
	if got := fromSlice([]int{}).Single(); got != 0 {
		t.Fatalf("empty Single() = %d, want zero value", got)
	}
	if got := fromSlice([]int{1, 2}).Single(); got != 0 {
		t.Fatalf("multiple Single() = %d, want zero value", got)
	}
}

func TestTypedSingleStopsAfterSecondElement(t *testing.T) {
	visited := 0
	q := Query[int]{iterate: func(yield func(int) bool) {
		for _, value := range []int{1, 2, 3} {
			visited++
			if !yield(value) {
				return
			}
		}
	}}

	if got := q.Single(); got != 0 {
		t.Fatalf("Single() = %d, want zero value", got)
	}
	if visited != 2 {
		t.Fatalf("source visited %d elements, want 2", visited)
	}
}

var benchmarkSingleResult int

func BenchmarkTypedSingle(b *testing.B) {
	source := make([]int, 1024)

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			result := FromSlice(source).Single()
			if result == nil {
				benchmarkSingleResult = 0
			}
		}
	})

	b.Run("generic_query", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkSingleResult = fromSlice(source).Single()
		}
	})
}
