package linq

import (
	"slices"
	"testing"
)

type selectManyPet struct{ name string }

type selectManyPerson struct {
	name string
	pets []selectManyPet
}

func (p selectManyPerson) Pets() Query[selectManyPet] { return fromSlice(p.pets) }

func petOwner(pet selectManyPet, person selectManyPerson) string {
	return person.name + ":" + pet.name
}

func TestTypedSelectManyByMethodExpression(t *testing.T) {
	source := []selectManyPerson{
		{name: "A", pets: []selectManyPet{{name: "one"}, {name: "two"}}},
		{name: "B", pets: []selectManyPet{{name: "three"}}},
	}

	got := fromSlice(source).
		SelectManyBy(selectManyPerson.Pets, petOwner).
		toSlice()
	want := []string{"A:one", "A:two", "B:three"}

	if !slices.Equal(got, want) {
		t.Fatalf("SelectManyBy() = %v, want %v", got, want)
	}
}

func TestTypedSelectManyByStopsBothIterators(t *testing.T) {
	selectorCalls := 0
	q := fromSlice([][]int{{1, 2}, {3, 4}}).SelectManyBy(
		func(values []int) Query[int] {
			selectorCalls++
			return fromSlice(values)
		},
		func(inner int, outer []int) int { return inner + len(outer) },
	)

	q.iterate(func(int) bool { return false })

	if selectorCalls != 1 {
		t.Fatalf("selector called %d times, want 1", selectorCalls)
	}
}

func selectManyResult(inner int, outer []int) int { return inner + len(outer) }

var benchmarkSelectManyBySum int

func BenchmarkTypedSelectManyBy(b *testing.B) {
	source := make([][]int, 128)
	for i := range source {
		source[i] = make([]int, 8)
		for j := range source[i] {
			source[i][j] = i*len(source[i]) + j
		}
	}

	b.Run("legacy_typed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			FromSlice(source).SelectManyByT(selectManyLegacy, selectManyResult).Iterate(func(value any) bool {
				sum += value.(int)
				return true
			})
			benchmarkSelectManyBySum = sum
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sum := 0
			fromSlice(source).SelectManyBy(selectManyTyped, selectManyResult).iterate(func(value int) bool {
				sum += value
				return true
			})
			benchmarkSelectManyBySum = sum
		}
	})
}
