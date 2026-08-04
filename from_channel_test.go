package linq

import (
	"slices"
	"testing"
)

func TestTypedFromChannelPreservesValues(t *testing.T) {
	source := make(chan string, 2)
	source <- "one"
	source <- "two"
	close(source)

	got := fromChannel(source).toSlice()
	if !slices.Equal(got, []string{"one", "two"}) {
		t.Fatalf("fromChannel() = %v, want [one two]", got)
	}
}

func TestTypedFromChannelLeavesValuesAfterConsumerStops(t *testing.T) {
	source := make(chan int, 2)
	source <- 1
	source <- 2
	close(source)

	fromChannel(source).iterate(func(int) bool { return false })
	if got := <-source; got != 2 {
		t.Fatalf("remaining channel value = %d, want 2", got)
	}
}

var benchmarkFromChannelSum int

func BenchmarkTypedFromChannel(b *testing.B) {
	const size = 256

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			source := make(chan int, size)
			for i := 0; i < size; i++ {
				source <- i
			}
			close(source)
			sum := 0
			legacyFromChannel(source).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkFromChannelSum = sum
		}
	})

	b.Run("generic_query", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			source := make(chan int, size)
			for i := 0; i < size; i++ {
				source <- i
			}
			close(source)
			sum := 0
			fromChannel(source).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkFromChannelSum = sum
		}
	})
}
