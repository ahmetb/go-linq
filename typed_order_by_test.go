package linq

import (
	"slices"
	"testing"
)

type orderWidget struct {
	rank int
	name string
}

func (w orderWidget) Rank() int { return w.rank }

func TestTypedOrderByMethodExpressionIsStable(t *testing.T) {
	source := []orderWidget{{rank: 2, name: "c"}, {rank: 1, name: "a"}, {rank: 1, name: "b"}}

	got := fromSlice(source).OrderBy(orderWidget.Rank).toSlice()
	want := []orderWidget{{rank: 1, name: "a"}, {rank: 1, name: "b"}, {rank: 2, name: "c"}}

	if !slices.Equal(got, want) {
		t.Fatalf("OrderBy(orderWidget.Rank) = %v, want %v", got, want)
	}
}

func TestTypedOrderByBuffersBeforeYielding(t *testing.T) {
	calls := 0
	q := fromSlice([]int{3, 2, 1}).OrderBy(func(value int) int {
		calls++
		return value
	})

	q.iterate(func(int) bool { return false })

	if calls == 0 {
		t.Fatal("selector was not called before ordered query yielded")
	}
}

func orderIdentity(value int) int { return value }

var benchmarkOrderBySum int

func BenchmarkTypedOrderBy(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = (i * 353) % len(source)
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			legacyFromSlice(source).OrderByT(orderIdentity).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkOrderBySum = sum
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).OrderBy(orderIdentity).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkOrderBySum = sum
		}
	})
}
