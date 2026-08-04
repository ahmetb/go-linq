package linq

import (
	"slices"
	"testing"
)

func TestTypedSkipWhileAcceptsMethodExpression(t *testing.T) {
	got := fromSlice([]allNumber{1, 2, -3, 4}).SkipWhile(allNumber.Positive).toSlice()
	want := []allNumber{-3, 4}

	if !slices.Equal(got, want) {
		t.Fatalf("SkipWhile(allNumber.Positive) = %v, want %v", got, want)
	}
}

func TestTypedSkipWhileStopsTestingAfterFirstFalse(t *testing.T) {
	calls := 0
	got := fromSlice([]int{1, 2, 3, 1}).SkipWhile(func(value int) bool {
		calls++
		return value < 3
	}).toSlice()

	want := []int{3, 1}
	if !slices.Equal(got, want) {
		t.Fatalf("SkipWhile() = %v, want %v", got, want)
	}
	if calls != 3 {
		t.Fatalf("predicate called %d times, want 3", calls)
	}
}

func TestTypedSkipWhileStopsWithConsumer(t *testing.T) {
	calls := 0
	q := fromSlice([]int{1, 2, 3}).SkipWhile(func(value int) bool {
		calls++
		return value < 2
	})

	q.iterate(func(int) bool { return false })

	if calls != 2 {
		t.Fatalf("predicate called %d times, want 2", calls)
	}
}

func skipWhileBenchmarkPredicate(value int) bool { return value < 512 }

var benchmarkSkipWhileSum int

func BenchmarkTypedSkipWhile(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			legacyFromSlice(source).SkipWhileT(skipWhileBenchmarkPredicate).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkSkipWhileSum = sum
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).SkipWhile(skipWhileBenchmarkPredicate).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkSkipWhileSum = sum
		}
	})
}
