package linq

import (
	"strconv"
	"testing"
)

func TestSelectMany(t *testing.T) {
	want := []int{1, 2, 3, 4, 5, 6, 7}
	if q := FromSlice([][]int{{1, 2, 3}, {4, 5, 6, 7}}).SelectMany(func(i []int) Query[int] {
		return FromSlice(i)
	}); !testQueryIteration(q, want) {
		t.Errorf("SelectMany()=%v expected %v", toSlice(q), want)
	}
}

func TestSelectMany_TypeChanging(t *testing.T) {
	want := []rune{'s', 't', 'r', 'i', 'n', 'g'}
	if q := FromSlice([]string{"str", "ing"}).SelectMany(func(s string) Query[rune] {
		return FromString(s)
	}); !testQueryIteration(q, want) {
		t.Errorf("SelectMany()=%v expected %v", toSlice(q), want)
	}
}

func TestSelectManyIndexed(t *testing.T) {
	tests := []struct {
		input    [][]int
		selector func(int, []int) Query[int]
		output   []int
	}{
		{[][]int{{1, 2, 3}, {4, 5, 6, 7}}, func(i int, x []int) Query[int] {
			if i > 0 {
				return FromSlice(x[1:])
			}
			return FromSlice(x)
		}, []int{1, 2, 3, 5, 6, 7}},
	}

	for _, test := range tests {
		if q := FromSlice(test.input).SelectManyIndexed(test.selector); !testQueryIteration(q, test.output) {
			t.Errorf("SelectManyIndexed()=%v expected %v", toSlice(q), test.output)
		}
	}

	want := []rune{'s', 't', 'r', '0', 'i', 'n', 'g', '1'}
	if q := FromSlice([]string{"str", "ing"}).SelectManyIndexed(func(i int, s string) Query[rune] {
		return FromString(s + strconv.Itoa(i))
	}); !testQueryIteration(q, want) {
		t.Errorf("SelectManyIndexed()=%v expected %v", toSlice(q), want)
	}
}

func TestSelectManyBy(t *testing.T) {
	want := []int{2, 3, 4, 5, 6, 7, 8}
	if q := FromSlice([][]int{{1, 2, 3}, {4, 5, 6, 7}}).SelectManyBy(
		func(i []int) Query[int] { return FromSlice(i) },
		func(x int, y []int) int { return x + 1 },
	); !testQueryIteration(q, want) {
		t.Errorf("SelectManyBy()=%v expected %v", toSlice(q), want)
	}
}

func TestSelectManyBy_TypeChanging(t *testing.T) {
	want := []string{"s_", "t_", "r_", "i_", "n_", "g_"}
	if q := FromSlice([]string{"str", "ing"}).SelectManyBy(
		func(s string) Query[rune] { return FromString(s) },
		func(r rune, s string) string { return string(r) + "_" },
	); !testQueryIteration(q, want) {
		t.Errorf("SelectManyBy()=%v expected %v", toSlice(q), want)
	}
}

func TestSelectManyByIndexed(t *testing.T) {
	want := []int{11, 21, 31, 5, 6, 7, 8}
	if q := FromSlice([][]int{{1, 2, 3}, {4, 5, 6, 7}}).SelectManyByIndexed(
		func(i int, x []int) Query[int] {
			if i == 0 {
				return FromSlice([]int{10, 20, 30})
			}
			return FromSlice(x)
		},
		func(x int, y []int) int { return x + 1 },
	); !testQueryIteration(q, want) {
		t.Errorf("SelectManyByIndexed()=%v expected %v", toSlice(q), want)
	}

	wantStr := []string{"s_", "t_", "r_", "i_", "n_", "g_"}
	if q := FromSlice([]string{"st", "ng"}).SelectManyByIndexed(
		func(i int, s string) Query[rune] {
			if i == 0 {
				return FromString(s + "r")
			}
			return FromString("i" + s)
		},
		func(r rune, s string) string { return string(r) + "_" },
	); !testQueryIteration(q, wantStr) {
		t.Errorf("SelectManyByIndexed()=%v expected %v", toSlice(q), wantStr)
	}
}
