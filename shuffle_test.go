package linq

import (
	"slices"
	"testing"
)

func TestShuffle(t *testing.T) {
	source := []int{1, 1, 2, 3, 5, 8}
	original := slices.Clone(source)
	q := FromSlice(source).Shuffle()

	orders := make([][]int, 20)
	for i := range orders {
		orders[i] = q.ToSlice()

		sorted := slices.Clone(orders[i])
		slices.Sort(sorted)
		if !slices.Equal(sorted, original) {
			t.Errorf("Shuffle()=%v expected permutation of %v", orders[i], original)
		}
	}
	// A Shuffle that does nothing, and one that shuffles once and caches the
	// result, both hold the order fixed across iterations. Two distinct orders
	// rule out both. With 360 distinguishable permutations of the source, a
	// false failure here runs at (1/360)^19.
	if !slices.ContainsFunc(orders, func(o []int) bool { return !slices.Equal(o, orders[0]) }) {
		t.Errorf("Shuffle yielded %v on all 20 iterations; expected the order to vary", orders[0])
	}
	if !slices.Equal(source, original) {
		t.Errorf("Shuffle modified source: got %v expected %v", source, original)
	}
	if got := FromSlice([]int{}).Shuffle().ToSlice(); got != nil {
		t.Errorf("Shuffle(empty)=%v expected nil", got)
	}
	if got := FromSlice([]int{42}).Shuffle().ToSlice(); !slices.Equal(got, []int{42}) {
		t.Errorf("Shuffle(single)=%v expected [42]", got)
	}
}

func TestShuffleBuffersBeforeYielding(t *testing.T) {
	pulled := 0
	q := FromSlice([]int{1, 2, 3, 4}).Where(func(int) bool {
		pulled++
		return true
	}).Shuffle()

	q.Iterate(func(int) bool { return false })
	if pulled != 4 {
		t.Errorf("Shuffle pulled %d elements expected 4", pulled)
	}
}

func TestShuffleTakeSamplesWholeSource(t *testing.T) {
	pulled := 0
	source := FromSeq(func(yield func(int) bool) {
		for i := range 100 {
			pulled++
			if !yield(i) {
				return
			}
		}
	})

	got := source.Shuffle().Take(5).ToSlice()
	if pulled != 100 {
		t.Errorf("Shuffle().Take(5) pulled %d elements expected 100", pulled)
	}
	if len(got) != 5 {
		t.Fatalf("Shuffle().Take(5) returned %d elements expected 5", len(got))
	}
	slices.Sort(got)
	for i, item := range got {
		if item < 0 || item >= 100 || i > 0 && item == got[i-1] {
			t.Fatalf("Shuffle().Take(5) returned invalid sample %v", got)
		}
	}
}
