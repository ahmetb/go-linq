package linq

import (
	"iter"
	"maps"
	"slices"
	"testing"
)

func TestPublicTypedConstructorsComposeWithGenericMethods(t *testing.T) {
	got := FromSlice([]int{1, 2, 3}).Select(func(value int) string {
		return string(rune('a' + value))
	}).Results()
	if !slices.Equal(got, []string{"b", "c", "d"}) {
		t.Fatalf("FromSlice().Select() = %v, want [b c d]", got)
	}

	mapResult := FromMap(map[int]string{1: "one"}).ToMapBy(
		func(item KeyValue[int, string]) int { return item.Key },
		func(item KeyValue[int, string]) string { return item.Value },
	)
	if !maps.Equal(mapResult, map[int]string{1: "one"}) {
		t.Fatalf("FromMap().ToMapBy() = %v", mapResult)
	}
}

func TestPublicFromSeqPreservesTypedSequence(t *testing.T) {
	sequence := iter.Seq[int](func(yield func(int) bool) {
		for _, value := range []int{1, 2, 3} {
			if !yield(value) {
				return
			}
		}
	})

	if got := FromSeq(sequence).Take(2).Results(); !slices.Equal(got, []int{1, 2}) {
		t.Fatalf("FromSeq().Take(2) = %v, want [1 2]", got)
	}
}

func TestPublicRangeAndRepeatReturnTypedQueries(t *testing.T) {
	got := Range(2, 2).Concat(Repeat(4, 2)).Results()
	if !slices.Equal(got, []int{2, 3, 4, 4}) {
		t.Fatalf("Range().Concat(Repeat()) = %v", got)
	}
}
