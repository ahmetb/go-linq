package linq

import "testing"

func TestTake(t *testing.T) {
	arr := [9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}

	tests := []struct {
		input  Query[int]
		output []int
	}{
		{FromSlice([]int{1, 2, 2, 3, 1}), []int{1, 2, 2}},
		{FromSlice(arr[:]), []int{1, 1, 1}},
	}

	for _, test := range tests {
		if q := test.input.Take(3); !testQueryIteration(q, test.output) {
			t.Errorf("Take(3)=%v expected %v", toSlice(q), test.output)
		}
	}

	want := []rune{'s', 's', 't'}
	if q := FromString("sstr").Take(3); !testQueryIteration(q, want) {
		t.Errorf("FromString(sstr).Take(3)=%v expected %v", toSlice(q), want)
	}
}

func TestTakeWhile(t *testing.T) {
	tests := []struct {
		input     Query[int]
		predicate func(int) bool
		output    []int
	}{
		{FromSlice([]int{1, 1, 1, 2, 1, 2}), func(i int) bool {
			return i < 3
		}, []int{1, 1, 1, 2, 1, 2}},
		{FromSlice([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}), func(i int) bool {
			return i < 3
		}, []int{1, 1, 1, 2, 1, 2}},
	}

	for _, test := range tests {
		if q := test.input.TakeWhile(test.predicate); !testQueryIteration(q, test.output) {
			t.Errorf("TakeWhile()=%v expected %v", toSlice(q), test.output)
		}
	}

	want := []rune{'s', 's'}
	if q := FromString("sstr").TakeWhile(func(r rune) bool {
		return r == 's'
	}); !testQueryIteration(q, want) {
		t.Errorf("FromString(sstr).TakeWhile()=%v expected %v", toSlice(q), want)
	}
}

func TestTakeWhileIndexed(t *testing.T) {
	tests := []struct {
		input     Query[int]
		predicate func(int, int) bool
		output    []int
	}{
		{FromSlice([]int{1, 1, 1, 2}), func(i int, x int) bool {
			return x < 2 || i < 5
		}, []int{1, 1, 1, 2}},
		{FromSlice([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}), func(i int, x int) bool {
			return x < 2 || i < 5
		}, []int{1, 1, 1, 2, 1}},
	}

	for _, test := range tests {
		if q := test.input.TakeWhileIndexed(test.predicate); !testQueryIteration(q, test.output) {
			t.Errorf("TakeWhileIndexed()=%v expected %v", toSlice(q), test.output)
		}
	}

	want := []rune{'s'}
	if q := FromString("sstr").TakeWhileIndexed(func(i int, r rune) bool {
		return r == 's' && i < 1
	}); !testQueryIteration(q, want) {
		t.Errorf("FromString(sstr).TakeWhileIndexed()=%v expected %v", toSlice(q), want)
	}
}
