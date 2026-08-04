package linq

import (
	"slices"
	"testing"
)

func TestTypedToChannelSendsValuesAndClosesChannel(t *testing.T) {
	result := make(chan int, 3)
	fromSlice([]int{3, 1, 2}).ToChannel(result)

	var got []int
	for {
		select {
		case value, open := <-result:
			if !open {
				if !slices.Equal(got, []int{3, 1, 2}) {
					t.Fatalf("ToChannel() sent %v, want [3 1 2]", got)
				}
				return
			}
			got = append(got, value)
		default:
			t.Fatalf("ToChannel() left channel open after sending %v", got)
		}
	}
}

func TestTypedToChannelAcceptsSendOnlyChannel(t *testing.T) {
	result := make(chan string, 1)
	var sendOnly chan<- string = result
	fromSlice([]string{"value"}).ToChannel(sendOnly)

	select {
	case got := <-result:
		if got != "value" {
			t.Fatalf("ToChannel() sent %q, want value", got)
		}
	default:
		t.Fatal("ToChannel() did not send a value")
	}
}

var benchmarkToChannelSum int

func BenchmarkTypedToChannel(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			result := make(chan int, len(source))
			FromSlice(source).ToChannelT(result)
			sum := 0
			for value := range result {
				sum += value
			}
			benchmarkToChannelSum = sum
		}
	})

	b.Run("generic_query", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			result := make(chan int, len(source))
			fromSlice(source).ToChannel(result)
			sum := 0
			for value := range result {
				sum += value
			}
			benchmarkToChannelSum = sum
		}
	})
}
