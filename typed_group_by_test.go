package linq

import (
	"slices"
	"strings"
	"testing"
)

func TestTypedGroupByPreservesEncounterOrder(t *testing.T) {
	got := fromSlice([]string{"one", "a", "two", "bb", "z"}).GroupBy(
		func(value string) int { return len(value) },
		strings.ToUpper,
	).toSlice()

	if len(got) != 3 {
		t.Fatalf("GroupBy() produced %d groups, want 3", len(got))
	}
	if got[0].Key != 3 || !slices.Equal(got[0].Group, []string{"ONE", "TWO"}) {
		t.Fatalf("first group = %#v, want key 3 and [ONE TWO]", got[0])
	}
	if got[1].Key != 1 || !slices.Equal(got[1].Group, []string{"A", "Z"}) {
		t.Fatalf("second group = %#v, want key 1 and [A Z]", got[1])
	}
	if got[2].Key != 2 || !slices.Equal(got[2].Group, []string{"BB"}) {
		t.Fatalf("third group = %#v, want key 2 and [BB]", got[2])
	}
}

func TestTypedGroupByBuffersBeforeYielding(t *testing.T) {
	calls := 0
	q := fromSlice([]int{1, 2, 3}).GroupBy(
		func(value int) int {
			calls++
			return value % 2
		},
		func(value int) int { return value },
	)

	q.iterate(func(group[int, int]) bool { return false })

	if calls != 3 {
		t.Fatalf("key selector called %d times, want 3", calls)
	}
}

func groupByBenchmarkKey(value int) int { return value % 128 }

func groupByBenchmarkElement(value int) int { return value * 2 }

var benchmarkGroupBySize int

func BenchmarkTypedGroupBy(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			size := 0
			FromSlice(source).GroupByT(groupByBenchmarkKey, groupByBenchmarkElement).Iterate(func(value any) bool {
				size += len(value.(Group).Group)
				return true
			})
			benchmarkGroupBySize = size
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			size := 0
			fromSlice(source).GroupBy(groupByBenchmarkKey, groupByBenchmarkElement).iterate(func(value group[int, int]) bool {
				size += len(value.Group)
				return true
			})
			benchmarkGroupBySize = size
		}
	})
}
