package linq

import (
	"strings"
	"testing"
)

func TestAggregate(t *testing.T) {
	input := []string{"apple", "mango", "orange", "passionfruit", "grape"}

	r, ok := FromSlice(input).Aggregate(func(r string, i string) string {
		if len(r) > len(i) {
			return r
		}
		return i
	})

	if !ok || r != "passionfruit" {
		t.Errorf("FromSlice(%v).Aggregate()=%v,%v expected passionfruit,true", input, r, ok)
	}
}

func TestAggregate_Empty(t *testing.T) {
	r, ok := FromSlice([]string{}).Aggregate(func(r string, i string) string {
		return r
	})

	if ok || r != "" {
		t.Errorf("FromSlice([]).Aggregate()=%v,%v expected \"\",false", r, ok)
	}
}

func TestAggregateWithSeed(t *testing.T) {
	input := []string{"apple", "mango", "orange", "banana", "grape"}
	want := "passionfruit"

	r := FromSlice(input).AggregateWithSeed(want,
		func(r string, i string) string {
			if len(r) > len(i) {
				return r
			}
			return i
		})

	if r != want {
		t.Errorf("FromSlice(%v).AggregateWithSeed()=%v expected %v", input, r, want)
	}
}

func TestAggregateWithSeed_TypeChanging(t *testing.T) {
	// The accumulator type (int) differs from the element type (string).
	input := []string{"apple", "mango", "orange"}
	want := 16

	r := FromSlice(input).AggregateWithSeed(0,
		func(acc int, i string) int {
			return acc + len(i)
		})

	if r != want {
		t.Errorf("FromSlice(%v).AggregateWithSeed()=%v expected %v", input, r, want)
	}
}

func TestAggregateWithSeedBy(t *testing.T) {
	input := []string{"apple", "mango", "orange", "passionfruit", "grape"}
	want := "PASSIONFRUIT"

	r := FromSlice(input).AggregateWithSeedBy("banana",
		func(r string, i string) string {
			if len(r) > len(i) {
				return r
			}
			return i
		},
		func(r string) string {
			return strings.ToUpper(r)
		},
	)

	if r != want {
		t.Errorf("FromSlice(%v).AggregateWithSeedBy()=%v expected %v", input, r, want)
	}
}
