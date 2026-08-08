package linq

import (
	"math"
	"testing"
)

func TestSum(t *testing.T) {
	tests := []struct {
		input []int
		want  int
	}{
		{[]int{1, 2, 2, 3, 1}, 9},
		{[]int{1}, 1},
		{[]int{}, 0},
	}

	for _, test := range tests {
		if r := Sum(FromSlice(test.input)); r != test.want {
			t.Errorf("Sum(FromSlice(%v))=%v expected %v", test.input, r, test.want)
		}
	}

	if r := Sum(FromSlice([]uint{1, 2, 2, 3, 1})); r != 9 {
		t.Errorf("Sum(uints)=%v expected 9", r)
	}

	if r := Sum(FromSlice([]float64{1., 2., 2., 3., 1.})); r != 9. {
		t.Errorf("Sum(floats)=%v expected 9.", r)
	}
}

func TestSumBy(t *testing.T) {
	input := []string{"apple", "banana", "fig"}
	want := 14

	if r := FromSlice(input).SumBy(func(s string) int {
		return len(s)
	}); r != want {
		t.Errorf("FromSlice(%v).SumBy()=%v expected %v", input, r, want)
	}
}

func TestAverage(t *testing.T) {
	if r := Average(FromSlice([]int{1, 2, 2, 3, 1})); r != 1.8 {
		t.Errorf("Average(ints)=%v expected 1.8", r)
	}

	if r := Average(FromSlice([]uint{1, 2, 5, 7, 10})); r != 5. {
		t.Errorf("Average(uints)=%v expected 5.", r)
	}

	if r := Average(FromSlice([]float32{1., 1.})); r != 1. {
		t.Errorf("Average(floats)=%v expected 1.", r)
	}
}

func TestAverage_LargeIntegerCancellation(t *testing.T) {
	const float64PrecisionLimit = int64(1 << 53)
	input := []int64{float64PrecisionLimit + 1, -float64PrecisionLimit}
	if got := Average(FromSlice(input)); got != 0.5 {
		t.Errorf("Average(%v)=%v; want 0.5", input, got)
	}
}

func TestAverage_LargeUnsignedPrecision(t *testing.T) {
	// Each element is 2^53+1, which is not representable in float64: converting
	// elements one by one silently drops the +1 three times, while summing in
	// uint64 first keeps the sum exact.
	v := uint64(1<<53) + 1
	input := []uint64{v, v, v}
	want := float64(v+v+v) / 3
	if got := Average(FromSlice(input)); got != want {
		t.Errorf("Average(%v)=%v; want %v", input, got, want)
	}
}

func TestAverageForNaN(t *testing.T) {
	if r := Average(FromSlice([]int{})); !math.IsNaN(r) {
		t.Errorf("Average(FromSlice([]int{}))=%v expected NaN", r)
	}
}

func TestAverageBy(t *testing.T) {
	input := []string{"apple", "banana", "fig"}
	want := 14. / 3.

	if r := FromSlice(input).AverageBy(func(s string) int {
		return len(s)
	}); r != want {
		t.Errorf("FromSlice(%v).AverageBy()=%v expected %v", input, r, want)
	}

	if r := FromSlice([]string{}).AverageBy(func(s string) int {
		return len(s)
	}); !math.IsNaN(r) {
		t.Errorf("AverageBy()=%v expected NaN", r)
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		input  []int
		want   int
		wantOk bool
	}{
		{[]int{1, 2, 2, 3, 1}, 3, true},
		{[]int{1}, 1, true},
		{[]int{}, 0, false},
	}

	for _, test := range tests {
		if r, ok := Max(FromSlice(test.input)); ok != test.wantOk || r != test.want {
			t.Errorf("Max(FromSlice(%v))=%v,%v expected %v,%v", test.input, r, ok, test.want, test.wantOk)
		}
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		input  []int
		want   int
		wantOk bool
	}{
		{[]int{1, 2, 2, 3, 0}, 0, true},
		{[]int{1}, 1, true},
		{[]int{}, 0, false},
	}

	for _, test := range tests {
		if r, ok := Min(FromSlice(test.input)); ok != test.wantOk || r != test.want {
			t.Errorf("Min(FromSlice(%v))=%v,%v expected %v,%v", test.input, r, ok, test.want, test.wantOk)
		}
	}
}

func TestMaxBy(t *testing.T) {
	input := []string{"apple", "banana", "fig"}

	if r, ok := FromSlice(input).MaxBy(func(s string) int {
		return len(s)
	}); !ok || r != "banana" {
		t.Errorf("MaxBy()=%v,%v expected banana,true", r, ok)
	}

	if r, ok := FromSlice([]string{}).MaxBy(func(s string) int {
		return len(s)
	}); ok || r != "" {
		t.Errorf("MaxBy()=%v,%v expected \"\",false", r, ok)
	}
}

func TestMinBy(t *testing.T) {
	input := []string{"apple", "banana", "fig"}

	if r, ok := FromSlice(input).MinBy(func(s string) int {
		return len(s)
	}); !ok || r != "fig" {
		t.Errorf("MinBy()=%v,%v expected fig,true", r, ok)
	}

	if r, ok := FromSlice([]string{}).MinBy(func(s string) int {
		return len(s)
	}); ok || r != "" {
		t.Errorf("MinBy()=%v,%v expected \"\",false", r, ok)
	}
}
