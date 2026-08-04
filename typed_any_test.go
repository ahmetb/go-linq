package linq

import "testing"

func TestTypedAnyReturnsTrueForNonEmptyQuery(t *testing.T) {
	if !fromSlice([]int{0}).Any() {
		t.Fatal("Any() = false, want true")
	}
}

func TestTypedAnyReturnsFalseForEmptyQuery(t *testing.T) {
	if fromSlice([]int{}).Any() {
		t.Fatal("Any() = true, want false")
	}
}

func TestTypedAnyStopsAfterFirstElement(t *testing.T) {
	visited := 0
	q := query[int]{iterate: func(yield func(int) bool) {
		for _, value := range []int{1, 2, 3} {
			visited++
			if !yield(value) {
				return
			}
		}
	}}

	if !q.Any() {
		t.Fatal("Any() = false, want true")
	}
	if visited != 1 {
		t.Fatalf("source visited %d elements, want 1", visited)
	}
}

var benchmarkAnyResult bool

func BenchmarkTypedAny(b *testing.B) {
	source := make([]int, 1024)

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkAnyResult = FromSlice(source).Any()
		}
	})

	b.Run("generic_query", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkAnyResult = fromSlice(source).Any()
		}
	})
}
