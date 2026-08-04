package linq

import "testing"

func TestWhere(t *testing.T) {
	arr := [9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}

	tests := []struct {
		input     Query[int]
		predicate func(int) bool
		output    []int
	}{
		{FromSlice(arr[:]), func(i int) bool {
			return i >= 3
		}, []int{3, 4}},
		{FromSlice([]int{1, 2, 3, 4}), func(i int) bool {
			return i%2 == 0
		}, []int{2, 4}},
	}

	for _, test := range tests {
		if q := test.input.Where(test.predicate); !testQueryIteration(q, test.output) {
			t.Errorf("Where()=%v expected %v", toSlice(q), test.output)
		}
	}
}

func TestWhere_String(t *testing.T) {
	want := []rune{'t', 'r'}
	if q := FromString("sstr").Where(func(r rune) bool {
		return r != 's'
	}); !testQueryIteration(q, want) {
		t.Errorf("FromString(sstr).Where()=%v expected %v", toSlice(q), want)
	}
}

func TestWhereIndexed(t *testing.T) {
	arr := [9]int{1, 1, 1, 2, 1, 2, 3, 4, 2}

	tests := []struct {
		input     Query[rune]
		predicate func(int, rune) bool
		output    []rune
	}{
		{FromString("sstr"), func(i int, x rune) bool {
			return x != 's' || i == 1
		}, []rune{'s', 't', 'r'}},
		{FromString("abcde"), func(i int, _ rune) bool {
			return i < 2
		}, []rune{'a', 'b'}},
	}

	for _, test := range tests {
		if q := test.input.WhereIndexed(test.predicate); !testQueryIteration(q, test.output) {
			t.Errorf("WhereIndexed()=%v expected %v", toSlice(q), test.output)
		}
	}

	want := []int{2, 3, 2}
	if q := FromSlice(arr[:]).WhereIndexed(func(i int, x int) bool {
		return x < 4 && i > 4
	}); !testQueryIteration(q, want) {
		t.Errorf("WhereIndexed()=%v expected %v", toSlice(q), want)
	}
}
