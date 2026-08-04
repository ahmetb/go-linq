package linq

import (
	"strconv"
	"testing"
)

func TestSelect(t *testing.T) {
	want := []int{2, 4, 6}
	if q := FromSlice([]int{1, 2, 3}).Select(func(i int) int {
		return i * 2
	}); !testQueryIteration(q, want) {
		t.Errorf("Select()=%v expected %v", toSlice(q), want)
	}
}

func TestSelect_TypeChanging(t *testing.T) {
	want := []string{"s1", "t1", "r1"}
	if q := FromString("str").Select(func(r rune) string {
		return string(r) + "1"
	}); !testQueryIteration(q, want) {
		t.Errorf("Select()=%v expected %v", toSlice(q), want)
	}
}

func TestSelectIndexed(t *testing.T) {
	want := []int{0, 2, 6}
	if q := FromSlice([]int{1, 2, 3}).SelectIndexed(func(i int, x int) int {
		return x * i
	}); !testQueryIteration(q, want) {
		t.Errorf("SelectIndexed()=%v expected %v", toSlice(q), want)
	}
}

func TestSelectIndexed_TypeChanging(t *testing.T) {
	want := []string{"s0", "t1", "r2"}
	if q := FromString("str").SelectIndexed(func(i int, r rune) string {
		return string(r) + strconv.Itoa(i)
	}); !testQueryIteration(q, want) {
		t.Errorf("SelectIndexed()=%v expected %v", toSlice(q), want)
	}
}
