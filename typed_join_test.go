package linq

import (
	"slices"
	"testing"
)

type joinPerson struct {
	id   int
	name string
}

func (p joinPerson) ID() int { return p.id }

type joinPet struct {
	ownerID int
	name    string
}

func (p joinPet) OwnerID() int { return p.ownerID }

func joinPersonPet(person joinPerson, pet joinPet) string { return person.name + ":" + pet.name }

func TestTypedJoinMethodExpressions(t *testing.T) {
	people := []joinPerson{{id: 2, name: "B"}, {id: 1, name: "A"}}
	pets := []joinPet{{ownerID: 1, name: "one"}, {ownerID: 2, name: "two"}, {ownerID: 2, name: "three"}}

	got := fromSlice(people).Join(
		fromSlice(pets),
		joinPerson.ID,
		joinPet.OwnerID,
		joinPersonPet,
	).toSlice()
	want := []string{"B:two", "B:three", "A:one"}

	if !slices.Equal(got, want) {
		t.Fatalf("Join() = %v, want %v", got, want)
	}
}

func TestTypedJoinStopsAfterFirstMatch(t *testing.T) {
	calls := 0
	q := fromSlice([]int{1, 2}).Join(
		fromSlice([]int{1, 1, 2}),
		func(value int) int { return value },
		func(value int) int { return value },
		func(left, right int) int {
			calls++
			return left + right
		},
	)

	q.iterate(func(int) bool { return false })

	if calls != 1 {
		t.Fatalf("result selector called %d times, want 1", calls)
	}
}

func joinBenchmarkKey(value int) int { return value % 128 }

func joinBenchmarkResult(outer, inner int) int { return outer + inner }

var benchmarkJoinSum int

func BenchmarkTypedJoin(b *testing.B) {
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
			legacyFromSlice(outer).JoinT(
				legacyFromSlice(inner), joinBenchmarkKey, joinBenchmarkKey, joinBenchmarkResult,
			).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkJoinSum = sum
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(outer).Join(
				fromSlice(inner), joinBenchmarkKey, joinBenchmarkKey, joinBenchmarkResult,
			).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkJoinSum = sum
		}
	})
}
