package linq

import (
	"context"
	"testing"
	"time"
)

func TestFromSlice(t *testing.T) {
	s := [3]int{1, 2, 3}
	w := []int{1, 2, 3}

	if q := FromSlice(s[:]); !testQueryIteration(q, w) {
		t.Errorf("FromSlice(%v)!=%v", s, w)
	}
}

func TestFromMap(t *testing.T) {
	s := map[string]bool{"foo": true}
	w := []KeyValue[string, bool]{{"foo", true}}

	if q := FromMap(s); !testQueryIteration(q, w) {
		t.Errorf("FromMap(%v)!=%v", s, w)
	}
}

func TestFromChannel(t *testing.T) {
	c := make(chan int, 3)
	c <- 10
	c <- 15
	c <- -3
	close(c)

	w := []int{10, 15, -3}

	if q := FromChannel(c); !assertQueryOutput(q, w) {
		t.Errorf("FromChannel() failed expected %v", w)
	}
}

func TestFromChannel_DryRun(t *testing.T) {
	c := make(chan int, 3)
	c <- 10
	c <- 15
	c <- -3
	close(c)
	q := FromChannel(c)
	runDryIteration(q)
}

func TestFromChannelWithContext_Cancel(t *testing.T) {
	c := make(chan int, 3)
	defer close(c)
	c <- 10
	c <- 15
	c <- -3

	w := []int{10, 15, -3}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if q := FromChannelWithContext(ctx, c); !assertQueryOutput(q, w) {
		t.Errorf("FromChannelWithContext() failed expected %v", w)
	}
}

func TestFromChannelWithContext_Closed(t *testing.T) {
	c := make(chan int, 3)
	c <- 10
	c <- 15
	c <- -3
	close(c)

	w := []int{10, 15, -3}

	ctx := context.Background()

	if q := FromChannelWithContext(ctx, c); !assertQueryOutput(q, w) {
		t.Errorf("FromChannelWithContext() failed expected %v", w)
	}
}

func TestFromString(t *testing.T) {
	s := "string"
	w := []rune{'s', 't', 'r', 'i', 'n', 'g'}

	if q := FromString(s); !testQueryIteration(q, w) {
		t.Errorf("FromString(%v)!=%v", s, w)
	}
}

func TestFromSeq(t *testing.T) {
	seq := func(yield func(int) bool) {
		for i := 1; i <= 3; i++ {
			if !yield(i) {
				return
			}
		}
	}
	w := []int{1, 2, 3}

	if q := FromSeq(seq); !testQueryIteration(q, w) {
		t.Errorf("FromSeq()!=%v", w)
	}
}

func TestRange(t *testing.T) {
	w := []int{-2, -1, 0, 1, 2}

	if q := Range(-2, 5); !testQueryIteration(q, w) {
		t.Errorf("Range(-2, 5)=%v expected %v", toSlice(q), w)
	}
}

func TestRepeat(t *testing.T) {
	w := []int{1, 1, 1, 1, 1}

	if q := Repeat(1, 5); !testQueryIteration(q, w) {
		t.Errorf("Repeat(1, 5)=%v expected %v", toSlice(q), w)
	}
}
