package linq

import (
	"fmt"
	"slices"
	"testing"
)

func groupJoinResult(person joinPerson, pets []joinPet) string {
	return fmt.Sprintf("%s:%d:%t", person.name, len(pets), pets != nil)
}

func TestTypedGroupJoinMethodExpressions(t *testing.T) {
	people := []joinPerson{{id: 1, name: "A"}, {id: 2, name: "B"}}
	pets := []joinPet{{ownerID: 1, name: "one"}, {ownerID: 1, name: "two"}}

	got := fromSlice(people).GroupJoin(
		fromSlice(pets), joinPerson.ID, joinPet.OwnerID, groupJoinResult,
	).toSlice()
	want := []string{"A:2:true", "B:0:true"}

	if !slices.Equal(got, want) {
		t.Fatalf("GroupJoin() = %v, want %v", got, want)
	}
}

func TestTypedGroupJoinStopsAfterFirstOuter(t *testing.T) {
	calls := 0
	q := fromSlice([]int{1, 2}).GroupJoin(
		fromSlice([]int{1, 1, 2}),
		func(value int) int { return value },
		func(value int) int { return value },
		func(outer int, inner []int) int {
			calls++
			return outer + len(inner)
		},
	)

	q.iterate(func(int) bool { return false })

	if calls != 1 {
		t.Fatalf("result selector called %d times, want 1", calls)
	}
}

func groupJoinBenchmarkResult(outer int, inner []int) int { return outer + len(inner) }

var benchmarkGroupJoinSum int

func BenchmarkTypedGroupJoin(b *testing.B) {
	outer := make([]int, 256)
	inner := make([]int, 1024)
	for i := range outer {
		outer[i] = i
	}
	for i := range inner {
		inner[i] = i
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			FromSlice(outer).GroupJoinT(
				FromSlice(inner), joinBenchmarkKey, joinBenchmarkKey, groupJoinBenchmarkResult,
			).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkGroupJoinSum = sum
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(outer).GroupJoin(
				fromSlice(inner), joinBenchmarkKey, joinBenchmarkKey, groupJoinBenchmarkResult,
			).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkGroupJoinSum = sum
		}
	})
}
