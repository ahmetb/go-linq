package linq

import (
	"slices"
	"testing"
)

func TestAppend(t *testing.T) {
	input := []int{1, 2, 3, 4}
	want := []int{1, 2, 3, 4, 5}

	if q := FromSlice(input).Append(5); !testQueryIteration(q, want) {
		t.Errorf("FromSlice(%v).Append()=%v expected %v", input, toSlice(q), want)
	}
}

func TestConcat(t *testing.T) {
	input1 := []int{1, 2, 3}
	input2 := []int{4, 5}
	want := []int{1, 2, 3, 4, 5}

	if q := FromSlice(input1).Concat(FromSlice(input2)); !testQueryIteration(q, want) {
		t.Errorf("FromSlice(%v).Concat(%v)=%v expected %v", input1, input2, toSlice(q), want)
	}
}

func TestPrepend(t *testing.T) {
	input := []int{1, 2, 3, 4}
	want := []int{0, 1, 2, 3, 4}

	if q := FromSlice(input).Prepend(0); !testQueryIteration(q, want) {
		t.Errorf("FromSlice(%v).Prepend()=%v expected %v", input, toSlice(q), want)
	}
}

// TestAppendEarlyExit verifies that stopping the consumer mid-source halts
// iteration immediately and the appended item is never yielded.
func TestAppendEarlyExit(t *testing.T) {
	var got []int
	FromSlice([]int{1, 2, 3}).Append(4).Iterate(func(v int) bool {
		got = append(got, v)
		return len(got) < 2
	})

	if want := []int{1, 2}; !slices.Equal(got, want) {
		t.Errorf("Append early exit yielded %v, expected %v", got, want)
	}
}

// TestConcatEarlyExit verifies that stopping the consumer during the first
// sequence halts iteration immediately and never starts the second sequence.
func TestConcatEarlyExit(t *testing.T) {
	firstYielded := 0
	first := Query[int]{Iterate: func(yield func(int) bool) {
		for i := 1; i <= 3; i++ {
			firstYielded++
			if !yield(i) {
				return
			}
		}
	}}
	secondStarted := false
	second := Query[int]{Iterate: func(yield func(int) bool) {
		secondStarted = true
		yield(4)
	}}

	var got []int
	first.Concat(second).Iterate(func(v int) bool {
		got = append(got, v)
		return len(got) < 2
	})

	if want := []int{1, 2}; !slices.Equal(got, want) {
		t.Errorf("Concat early exit yielded %v, expected %v", got, want)
	}
	if firstYielded != 2 {
		t.Errorf("first sequence yielded %d items after consumer stop, expected 2", firstYielded)
	}
	if secondStarted {
		t.Error("second sequence was started despite early exit in the first")
	}
}

// TestPrependEarlyExit verifies that stopping the consumer on the prepended
// item never starts the source, and that stopping mid-source halts iteration
// immediately.
func TestPrependEarlyExit(t *testing.T) {
	sourceStarted := false
	source := Query[int]{Iterate: func(yield func(int) bool) {
		sourceStarted = true
		for i := 1; i <= 3; i++ {
			if !yield(i) {
				return
			}
		}
	}}

	var got []int
	source.Prepend(0).Iterate(func(v int) bool {
		got = append(got, v)
		return false
	})

	if want := []int{0}; !slices.Equal(got, want) {
		t.Errorf("Prepend early exit yielded %v, expected %v", got, want)
	}
	if sourceStarted {
		t.Error("source was started despite early exit on the prepended item")
	}

	got = nil
	source.Prepend(0).Iterate(func(v int) bool {
		got = append(got, v)
		return len(got) < 2
	})

	if want := []int{0, 1}; !slices.Equal(got, want) {
		t.Errorf("Prepend mid-source early exit yielded %v, expected %v", got, want)
	}
}
