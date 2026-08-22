package linq

import (
	"slices"
	"testing"
)

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

func TestTakeLast(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	tests := []struct {
		count int
		want  []int
	}{
		{-1, nil},
		{0, nil},
		{1, []int{5}},
		{3, []int{3, 4, 5}},
		{5, []int{1, 2, 3, 4, 5}},
		{10, []int{1, 2, 3, 4, 5}},
	}

	for _, test := range tests {
		if q := FromSlice(input).TakeLast(test.count); !testQueryIteration(q, test.want) {
			t.Errorf("TakeLast(%d)=%v expected %v", test.count, toSlice(q), test.want)
		}
	}
}

func TestTakeLastConsumesSourceBeforeYielding(t *testing.T) {
	pulled := 0
	q := FromSlice([]int{1, 2, 3, 4}).Where(func(int) bool {
		pulled++
		return true
	}).TakeLast(2)

	var got []int
	q.Iterate(func(item int) bool {
		got = append(got, item)
		return false
	})

	if !slices.Equal(got, []int{3}) {
		t.Errorf("TakeLast early exit yielded %v expected [3]", got)
	}
	if pulled != 4 {
		t.Errorf("source pulled %d elements expected 4", pulled)
	}
}

func TestTakeRange(t *testing.T) {
	input := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	tests := []struct {
		start Position
		end   Position
		want  []int
	}{
		{FromStart(2), FromStart(7), []int{2, 3, 4, 5, 6}},
		{FromStart(2), FromEnd(3), []int{2, 3, 4, 5, 6}},
		{FromEnd(7), FromEnd(3), []int{3, 4, 5, 6}},
		{FromEnd(7), FromStart(8), []int{3, 4, 5, 6, 7}},
		{FromStart(0), FromEnd(0), input},
		{FromStart(20), FromEnd(0), nil},
		{FromEnd(20), FromEnd(0), input},
		{FromEnd(2), FromEnd(5), nil},
		{FromStart(7), FromStart(2), nil},
		{FromStart(2), FromStart(20), []int{2, 3, 4, 5, 6, 7, 8, 9}},
		{FromStart(2), FromEnd(20), nil},
		{FromEnd(20), FromStart(3), []int{0, 1, 2}},
	}

	sources := []Query[int]{
		FromSlice(input),
		FromSeq(FromSlice(input).Iterate),
	}
	for _, source := range sources {
		for _, test := range tests {
			q := source.TakeRange(test.start, test.end)
			if !testQueryIteration(q, test.want) {
				t.Errorf("TakeRange(%v, %v)=%v expected %v",
					test.start, test.end,
					toSlice(q), test.want)
			}
		}
	}
}

func TestTakeRangeFromStartStopsAtEnd(t *testing.T) {
	pulled := 0
	q := FromSlice([]int{0, 1, 2, 3, 4, 5}).Where(func(int) bool {
		pulled++
		return true
	}).TakeRange(FromStart(2), FromStart(5))

	if got := q.ToSlice(); !slices.Equal(got, []int{2, 3, 4}) {
		t.Errorf("TakeRange(FromStart(2), FromStart(5))=%v expected [2 3 4]", got)
	}
	if pulled != 5 {
		t.Errorf("source pulled %d elements expected 5", pulled)
	}
}

// TestTakePullsExactlyCount verifies Take stops pulling from the source once
// it has yielded count elements, rather than pulling one extra element just to
// discard it. The Where predicate counts how many elements the source produced.
func TestTakePullsExactlyCount(t *testing.T) {
	tests := []struct {
		count  int
		output []int
		pulled int
	}{
		{0, nil, 0},
		{1, []int{1}, 1},
		{3, []int{1, 2, 3}, 3},
		{4, []int{1, 2, 3, 4}, 4},
		{9, []int{1, 2, 3, 4}, 4},
	}

	for _, test := range tests {
		pulled := 0
		q := FromSlice([]int{1, 2, 3, 4}).Where(func(int) bool {
			pulled++
			return true
		}).Take(test.count)

		if out := toSlice(q); !slices.Equal(out, test.output) {
			t.Errorf("Take(%d)=%v expected %v", test.count, out, test.output)
		}
		if pulled != test.pulled {
			t.Errorf("Take(%d) pulled %d elements from the source, expected %d",
				test.count, pulled, test.pulled)
		}
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
