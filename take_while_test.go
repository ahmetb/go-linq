package linq

import (
	"slices"
	"testing"
)

func TestTypedTakeWhileAcceptsMethodExpression(t *testing.T) {
	got := fromSlice([]allNumber{1, 2, -3, 4}).TakeWhile(allNumber.Positive).toSlice()
	want := []allNumber{1, 2}

	if !slices.Equal(got, want) {
		t.Fatalf("TakeWhile(allNumber.Positive) = %v, want %v", got, want)
	}
}

func TestTypedTakeWhileStopsAtFirstFalse(t *testing.T) {
	calls := 0
	got := fromSlice([]int{1, 2, 3, 1}).TakeWhile(func(value int) bool {
		calls++
		return value < 3
	}).toSlice()

	want := []int{1, 2}
	if !slices.Equal(got, want) {
		t.Fatalf("TakeWhile() = %v, want %v", got, want)
	}
	if calls != 3 {
		t.Fatalf("predicate called %d times, want 3", calls)
	}
}

func TestTypedTakeWhileStopsWithConsumer(t *testing.T) {
	calls := 0
	q := fromSlice([]int{1, 2, 3}).TakeWhile(func(value int) bool {
		calls++
		return true
	})

	q.iterate(func(int) bool { return false })

	if calls != 1 {
		t.Fatalf("predicate called %d times, want 1", calls)
	}
}

func takeWhileBenchmarkPredicate(value int) bool { return value < 512 }

var benchmarkTakeWhileSum int

func BenchmarkTypedTakeWhile(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			legacyFromSlice(source).TakeWhileT(takeWhileBenchmarkPredicate).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkTakeWhileSum = sum
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).TakeWhile(takeWhileBenchmarkPredicate).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkTakeWhileSum = sum
		}
	})
}
