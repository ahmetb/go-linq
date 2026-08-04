package linq

import (
	"slices"
	"testing"
)

func TestTypedFromSliceToSlice(t *testing.T) {
	type numbers []int

	source := numbers{1, 2, 3}
	got := fromSlice(source).toSlice()
	want := []int{1, 2, 3}

	if !slices.Equal(got, want) {
		t.Fatalf("fromSlice(%v).toSlice() = %v, want %v", source, got, want)
	}
}

func TestTypedFromSliceStopsWhenConsumerStops(t *testing.T) {
	q := fromSlice([]int{1, 2, 3})
	var got []int

	q.iterate(func(value int) bool {
		got = append(got, value)
		return false
	})

	want := []int{1}
	if !slices.Equal(got, want) {
		t.Fatalf("iteration produced %v, want %v", got, want)
	}
}
