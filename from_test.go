package linq

import (
	"context"
	"testing"
	"time"
)

func TestFromSlice(t *testing.T) {
	s := [3]int{1, 2, 3}
	w := []any{1, 2, 3}

	if q := FromSlice(s[:]); !validateQuery(q, w) {
		t.Errorf("FromSlice(%v)!=%v", s, w)
	}
}

func TestFromMap(t *testing.T) {
	s := map[string]bool{"foo": true}
	w := []any{KeyValue{"foo", true}}

	if q := FromMap(s); !validateQuery(q, w) {
		t.Errorf("FromMap(%v)!=%v", s, w)
	}
}

func TestFromChannel(t *testing.T) {
	c := make(chan int, 3)
	c <- 10
	c <- 15
	c <- -3
	close(c)

	w := []any{10, 15, -3}

	if q := FromChannel(c); !validateQuery(q, w) {
		t.Errorf("FromChannel() failed expected %v", w)
	}
}

func TestFromChannelWithContext_Cancel(t *testing.T) {
	c := make(chan int, 3)
	defer close(c)
	c <- 10
	c <- 15
	c <- -3

	w := []any{10, 15, -3}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if q := FromChannelWithContext(c, ctx); !validateQuery(q, w) {
		t.Errorf("FromChannelWithContext() failed expected %v", w)
	}
}

func TestFromChannelWithContext_Closed(t *testing.T) {
	c := make(chan int, 3)
	c <- 10
	c <- 15
	c <- -3
	close(c)

	w := []any{10, 15, -3}

	ctx := context.Background()

	if q := FromChannelWithContext(c, ctx); !validateQuery(q, w) {
		t.Errorf("FromChannelWithContext() failed expected %v", w)
	}
}

func TestFromChannelWithTimeout_Timeout(t *testing.T) {
	c := make(chan int, 3)
	defer close(c)
	c <- 10
	c <- 15
	c <- -3

	w := []any{10, 15, -3}

	if q := FromChannelWithTimeout(c, time.Second); !validateQuery(q, w) {
		t.Errorf("FromChannelWithTimeout() failed expected %v", w)
	}
}

func TestFromChannelWithTimeout_Closed(t *testing.T) {
	c := make(chan int, 3)
	c <- 10
	c <- 15
	c <- -3
	close(c)

	w := []any{10, 15, -3}

	if q := FromChannelWithTimeout(c, time.Hour); !validateQuery(q, w) {
		t.Errorf("FromChannelWithTimeout() failed expected %v", w)
	}
}

func TestFromString(t *testing.T) {
	s := "string"
	w := []any{'s', 't', 'r', 'i', 'n', 'g'}

	if q := FromString(s); !validateQuery(q, w) {
		t.Errorf("FromString(%v)!=%v", s, w)
	}
}

func TestFromIterable(t *testing.T) {
	s := foo{f1: 1, f2: true, f3: "string"}
	w := []any{1, true, "string"}

	if q := FromIterable(s); !validateQuery(q, w) {
		t.Errorf("FromIterable(%v)!=%v", s, w)
	}
}

func TestFrom(t *testing.T) {
	c := make(chan any, 3)
	c <- -1
	c <- 0
	c <- 1
	close(c)

	ct := make(chan int, 3)
	ct <- -10
	ct <- 0
	ct <- 10
	close(ct)

	tests := []struct {
		input  any
		output []any
		want   bool
	}{
		{[]int{1, 2, 3}, []any{1, 2, 3}, true},
		{[]int{1, 2, 4}, []any{1, 2, 3}, false},
		{[3]int{1, 2, 3}, []any{1, 2, 3}, true},
		{[3]int{1, 2, 4}, []any{1, 2, 3}, false},
		{"str", []any{'s', 't', 'r'}, true},
		{"str", []any{'s', 't', 'g'}, false},
		{map[string]bool{"foo": true}, []any{KeyValue{"foo", true}}, true},
		{map[string]bool{"foo": true}, []any{KeyValue{"foo", false}}, false},
		{c, []any{-1, 0, 1}, true},
		{ct, []any{-10, 0, 10}, true},
		{foo{f1: 1, f2: true, f3: "string"}, []any{1, true, "string"}, true},
	}

	for _, test := range tests {
		if q := From(test.input); validateQuery(q, test.output) != test.want {
			if test.want {
				t.Errorf("From(%v)=%v expected %v", test.input, toSlice(q), test.output)
			} else {
				t.Errorf("From(%v)=%v expected not equal", test.input, test.output)
			}
		}
	}
}

func TestRange(t *testing.T) {
	w := []any{-2, -1, 0, 1, 2}

	if q := Range(-2, 5); !validateQuery(q, w) {
		t.Errorf("Range(-2, 5)=%v expected %v", toSlice(q), w)
	}
}

func TestRepeat(t *testing.T) {
	w := []any{1, 1, 1, 1, 1}

	if q := Repeat(1, 5); !validateQuery(q, w) {
		t.Errorf("Repeat(1, 5)=%v expected %v", toSlice(q), w)
	}
}
