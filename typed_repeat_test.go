package linq

import (
	"slices"
	"testing"
)

func TestTypedRepeatYieldsRequestedCopies(t *testing.T) {
	got := repeatValue("value", 3).toSlice()
	if !slices.Equal(got, []string{"value", "value", "value"}) {
		t.Fatalf("repeatValue() = %v, want three values", got)
	}
}

func TestTypedRepeatNonPositiveCountIsEmpty(t *testing.T) {
	if got := repeatValue(1, -1).toSlice(); len(got) != 0 {
		t.Fatalf("repeatValue(1, -1) = %v, want empty", got)
	}
}

func TestTypedRepeatStopsWithConsumer(t *testing.T) {
	calls := 0
	repeatValue([]int{1}, 3).iterate(func([]int) bool {
		calls++
		return false
	})

	if calls != 1 {
		t.Fatalf("consumer called %d times, want 1", calls)
	}
}

var benchmarkRepeatSum int

func BenchmarkTypedRepeat(b *testing.B) {
	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			legacyRepeat(1, 1024).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkRepeatSum = sum
		}
	})

	b.Run("generic_query", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			repeatValue(1, 1024).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkRepeatSum = sum
		}
	})
}
