package linq

import (
	"context"
	"slices"
	"testing"
)

func TestTypedFromChannelWithContextPreservesValues(t *testing.T) {
	source := make(chan int, 2)
	source <- 1
	source <- 2
	close(source)

	got := fromChannelWithContext(context.Background(), source).toSlice()
	if !slices.Equal(got, []int{1, 2}) {
		t.Fatalf("fromChannelWithContext() = %v, want [1 2]", got)
	}
}

func TestTypedFromChannelWithContextStopsWhenCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	source := make(chan int)

	if got := fromChannelWithContext(ctx, source).toSlice(); len(got) != 0 {
		t.Fatalf("canceled query = %v, want empty", got)
	}
}

func TestTypedFromChannelWithContextStopsWithConsumer(t *testing.T) {
	source := make(chan string, 2)
	source <- "one"
	source <- "two"
	close(source)

	fromChannelWithContext(context.Background(), source).iterate(func(string) bool { return false })
	if got := <-source; got != "two" {
		t.Fatalf("remaining channel value = %q, want two", got)
	}
}

var benchmarkFromChannelContextSum int

func BenchmarkTypedFromChannelWithContext(b *testing.B) {
	const size = 256
	ctx := context.Background()

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			source := make(chan int, size)
			for i := 0; i < size; i++ {
				source <- i
			}
			close(source)
			sum := 0
			legacyFromChannelWithContext(ctx, source).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkFromChannelContextSum = sum
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
			fromChannelWithContext(ctx, source).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkFromChannelContextSum = sum
		}
	})
}
