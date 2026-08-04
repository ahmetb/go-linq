package linq

import "testing"

func TestZip(t *testing.T) {
	input1 := []int{1, 2, 3}
	input2 := []int{2, 4, 5, 1}
	want := []int{3, 6, 8}

	if q := FromSlice(input1).Zip(FromSlice(input2), func(i, j int) int {
		return i + j
	}); !testQueryIteration(q, want) {
		t.Errorf("FromSlice(%v).Zip(%v)=%v expected %v", input1, input2, toSlice(q), want)
	}
}

func TestZip_TypeChanging(t *testing.T) {
	input1 := []string{"a", "b", "c"}
	input2 := []int{1, 2, 3}
	want := []string{"a1", "b2", "c3"}

	if q := FromSlice(input1).Zip(FromSlice(input2), func(s string, i int) string {
		return s + string(rune('0'+i))
	}); !testQueryIteration(q, want) {
		t.Errorf("Zip()=%v expected %v", toSlice(q), want)
	}
}
