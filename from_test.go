package linq

import (
	"context"
	"math"
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

func TestEmptyQuery(t *testing.T) {
	q := Empty[int]()
	if !testQueryIteration(q, nil) {
		t.Error("Empty[int]() yielded elements")
	}
	if q.Any() || q.Count() != 0 {
		t.Error("Empty[int]() is not empty")
	}
	if appended := q.Append(1); !testQueryIteration(appended, []int{1}) {
		t.Errorf("Empty[int]().Append(1)=%v expected [1]", appended.ToSlice())
	}
}

func TestRange(t *testing.T) {
	w := []int{-2, -1, 0, 1, 2}

	if q := Range(-2, 5); !testQueryIteration(q, w) {
		t.Errorf("Range(-2, 5)=%v expected %v", toSlice(q), w)
	}
}

func TestSequence(t *testing.T) {
	tests := []struct {
		name             string
		start, end, step int
		want             []int
	}{
		{"increasing", 1, 5, 1, []int{1, 2, 3, 4, 5}},
		{"decreasing", 5, 1, -2, []int{5, 3, 1}},
		{"bound not reached", 0, 5, 2, []int{0, 2, 4}},
		{"equal positive", 3, 3, 1, []int{3}},
		{"equal negative", 3, 3, -1, []int{3}},
		{"equal zero", 3, 3, 0, []int{3}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			q := Sequence(test.start, test.end, test.step)
			if !testQueryIteration(q, test.want) {
				t.Errorf("Sequence(%v, %v, %v)=%v expected %v",
					test.start, test.end, test.step, q.ToSlice(), test.want)
			}
		})
	}
}

func TestSequenceNumericTypes(t *testing.T) {
	type score int16

	if q := Sequence(score(-2), score(4), score(3)); !testQueryIteration(q, []score{-2, 1, 4}) {
		t.Errorf("Sequence(named type)=%v expected [-2 1 4]", q.ToSlice())
	}
	if q := Sequence(0.5, 1.25, 0.25); !testQueryIteration(q, []float64{0.5, 0.75, 1, 1.25}) {
		t.Errorf("Sequence(float64)=%v expected [0.5 0.75 1 1.25]", q.ToSlice())
	}
	if q := Sequence(uint8(254), uint8(255), uint8(2)); !testQueryIteration(q, []uint8{254}) {
		t.Errorf("Sequence(uint8 overflow)=%v expected [254]", q.ToSlice())
	}
	if q := Sequence(int8(126), int8(127), int8(2)); !testQueryIteration(q, []int8{126}) {
		t.Errorf("Sequence(int8 overflow)=%v expected [126]", q.ToSlice())
	}
}

func TestSequenceInvalidArgumentsPanic(t *testing.T) {
	for _, f := range []func(){
		func() { Sequence(0, 1, 0) },
		func() { Sequence(1, 0, 1) },
		func() { Sequence(0, 1, -1) },
		func() { Sequence(math.NaN(), 1, 1) },
		func() { Sequence(0, math.NaN(), 1) },
		func() { Sequence(0, 1, math.NaN()) },
	} {
		mustPanic(t, f)
	}
}

func TestInfiniteSequence(t *testing.T) {
	tests := []struct {
		name        string
		start, step int
		count       int
		want        []int
	}{
		{"increasing", 1, 2, 5, []int{1, 3, 5, 7, 9}},
		{"decreasing", 5, -2, 4, []int{5, 3, 1, -1}},
		{"constant", 7, 0, 3, []int{7, 7, 7}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			q := InfiniteSequence(test.start, test.step).Take(test.count)
			if !testQueryIteration(q, test.want) {
				t.Errorf("InfiniteSequence(%v, %v).Take(%v)=%v expected %v",
					test.start, test.step, test.count, q.ToSlice(), test.want)
			}
		})
	}

	q := InfiniteSequence(uint8(254), uint8(2)).Take(3)
	if !testQueryIteration(q, []uint8{254, 0, 2}) {
		t.Errorf("InfiniteSequence(uint8 overflow)=%v expected [254 0 2]", q.ToSlice())
	}
}

func TestRepeat(t *testing.T) {
	w := []int{1, 1, 1, 1, 1}

	if q := Repeat(1, 5); !testQueryIteration(q, w) {
		t.Errorf("Repeat(1, 5)=%v expected %v", toSlice(q), w)
	}
}
