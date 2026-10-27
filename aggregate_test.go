package linq

import (
	"slices"
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

func TestCountBy(t *testing.T) {
	input := []string{"apple", "banana", "apricot", "blueberry", "avocado", "cherry"}
	want := []KeyValue[byte, int]{
		{Key: 'a', Value: 3},
		{Key: 'b', Value: 2},
		{Key: 'c', Value: 1},
	}

	q := FromSlice(input).CountBy(func(item string) byte { return item[0] })
	if !testQueryIteration(q, want) {
		t.Errorf("CountBy()=%v expected %v", q.ToSlice(), want)
	}
	if got := FromSlice([]int{}).CountBy(func(item int) int { return item }).ToSlice(); got != nil {
		t.Errorf("CountBy(empty)=%v expected nil", got)
	}
}

func TestAggregateBy(t *testing.T) {
	type score struct {
		id    string
		value int
	}
	input := []score{{"0", 42}, {"1", 5}, {"2", 4}, {"1", 10}, {"0", 25}}
	want := []KeyValue[string, int]{
		{Key: "0", Value: 67},
		{Key: "1", Value: 15},
		{Key: "2", Value: 4},
	}

	q := FromSlice(input).AggregateBy(
		func(item score) string { return item.id },
		0,
		func(total int, item score) int { return total + item.value },
	)
	if !testQueryIteration(q, want) {
		t.Errorf("AggregateBy()=%v expected %v", q.ToSlice(), want)
	}
	if got := FromSlice([]int{}).AggregateBy(
		func(item int) int { return item },
		0,
		func(total, item int) int { return total + item },
	).ToSlice(); got != nil {
		t.Errorf("AggregateBy(empty)=%v expected nil", got)
	}
}

func TestAggregateByWithSeedSelector(t *testing.T) {
	type item struct {
		key   string
		value int
	}
	seedCalls := make(map[string]int)
	q := FromSlice([]item{{"a", 1}, {"bb", 2}, {"a", 3}}).AggregateByWithSeedSelector(
		func(item item) string { return item.key },
		func(key string) int {
			seedCalls[key]++
			return len(key)
		},
		func(total int, item item) int { return total + item.value },
	)

	want := []KeyValue[string, int]{{Key: "a", Value: 5}, {Key: "bb", Value: 4}}
	if got := q.ToSlice(); !slices.Equal(got, want) {
		t.Errorf("AggregateByWithSeedSelector()=%v expected %v", got, want)
	}
	for _, key := range []string{"a", "bb"} {
		if seedCalls[key] != 1 {
			t.Errorf("seedSelector called %d times for %q expected 1", seedCalls[key], key)
		}
	}
}

func TestCountByConsumesSourceBeforeYielding(t *testing.T) {
	pulled := 0
	q := FromSlice([]int{1, 2, 3, 4}).Where(func(int) bool {
		pulled++
		return true
	}).CountBy(func(item int) int { return item % 2 })

	q.Iterate(func(KeyValue[int, int]) bool { return false })
	if pulled != 4 {
		t.Errorf("CountBy pulled %d elements expected 4", pulled)
	}
}

func TestAggregateByConsumesSourceBeforeYielding(t *testing.T) {
	pulled := 0
	q := FromSlice([]int{1, 2, 3, 4}).Where(func(int) bool {
		pulled++
		return true
	}).AggregateBy(
		func(item int) int { return item % 2 },
		0,
		func(total, item int) int { return total + item },
	)

	q.Iterate(func(KeyValue[int, int]) bool { return false })
	if pulled != 4 {
		t.Errorf("AggregateBy pulled %d elements expected 4", pulled)
	}
}
