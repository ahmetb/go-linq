package linq

import "testing"

func TestSkip(t *testing.T) {
	arr := [9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}

	tests := []struct {
		input  Query[int]
		output []int
	}{
		{FromSlice([]int{1, 2}), []int{}},
		{FromSlice([]int{1, 2, 2, 3, 1}), []int{3, 1}},
		{FromSlice(arr[:]), []int{2, 1, 2, 3, 4, 2}},
	}

	for _, test := range tests {
		if q := test.input.Skip(3); !testQueryIteration(q, test.output) {
			t.Errorf("Skip(3)=%v expected %v", toSlice(q), test.output)
		}
	}

	want := []rune{'r'}
	if q := FromString("sstr").Skip(3); !testQueryIteration(q, want) {
		t.Errorf("FromString(sstr).Skip(3)=%v expected %v", toSlice(q), want)
	}
}

func TestSkipLast(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	tests := []struct {
		count int
		want  []int
	}{
		{-1, []int{1, 2, 3, 4, 5}},
		{0, []int{1, 2, 3, 4, 5}},
		{1, []int{1, 2, 3, 4}},
		{3, []int{1, 2}},
		{5, nil},
		{10, nil},
	}

	for _, test := range tests {
		if q := FromSlice(input).SkipLast(test.count); !testQueryIteration(q, test.want) {
			t.Errorf("SkipLast(%d)=%v expected %v", test.count, toSlice(q), test.want)
		}
	}
}

func TestSkipLastStopsPullingAfterConsumerStops(t *testing.T) {
	pulled := 0
	q := FromSlice([]int{1, 2, 3, 4, 5}).Where(func(int) bool {
		pulled++
		return true
	}).SkipLast(2)

	var got []int
	q.Iterate(func(item int) bool {
		got = append(got, item)
		return false
	})

	if len(got) != 1 || got[0] != 1 {
		t.Errorf("SkipLast early exit yielded %v expected [1]", got)
	}
	if pulled != 3 {
		t.Errorf("source pulled %d elements expected 3", pulled)
	}
}

func TestSkipWhile(t *testing.T) {
	tests := []struct {
		input     Query[int]
		predicate func(int) bool
		output    []int
	}{
		{FromSlice([]int{1, 2}), func(i int) bool {
			return i < 3
		}, []int{}},
		{FromSlice([]int{4, 1, 2}), func(i int) bool {
			return i < 3
		}, []int{4, 1, 2}},
		{FromSlice([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}), func(i int) bool {
			return i < 3
		}, []int{3, 4, 2}},
	}

	for _, test := range tests {
		if q := test.input.SkipWhile(test.predicate); !testQueryIteration(q, test.output) {
			t.Errorf("SkipWhile()=%v expected %v", toSlice(q), test.output)
		}
	}

	want := []rune{'t', 'r'}
	if q := FromString("sstr").SkipWhile(func(r rune) bool {
		return r == 's'
	}); !testQueryIteration(q, want) {
		t.Errorf("FromString(sstr).SkipWhile()=%v expected %v", toSlice(q), want)
	}
}

func TestSkipWhileIndexed(t *testing.T) {
	tests := []struct {
		input     Query[int]
		predicate func(int, int) bool
		output    []int
	}{
		{FromSlice([]int{1, 2}), func(i int, x int) bool {
			return x < 3
		}, []int{}},
		{FromSlice([]int{4, 1, 2}), func(i int, x int) bool {
			return x < 3
		}, []int{4, 1, 2}},
		{FromSlice([]int{1, 1, 1, 2, 1, 2, 3, 4, 2}), func(i int, x int) bool {
			return x < 2 || i < 5
		}, []int{2, 3, 4, 2}},
	}

	for _, test := range tests {
		if q := test.input.SkipWhileIndexed(test.predicate); !testQueryIteration(q, test.output) {
			t.Errorf("SkipWhileIndexed()=%v expected %v", toSlice(q), test.output)
		}
	}

	want := []rune{'s', 't', 'r'}
	if q := FromString("sstr").SkipWhileIndexed(func(i int, r rune) bool {
		return r == 's' && i < 1
	}); !testQueryIteration(q, want) {
		t.Errorf("FromString(sstr).SkipWhileIndexed()=%v expected %v", toSlice(q), want)
	}
}
