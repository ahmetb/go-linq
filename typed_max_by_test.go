package linq

import "testing"

func TestTypedMaxByAcceptsMethodExpression(t *testing.T) {
	people := []joinPerson{{id: 2, name: "two"}, {id: 5, name: "five"}, {id: 3, name: "three"}}
	got := fromSlice(people).MaxBy(joinPerson.ID)

	if got != people[1] {
		t.Fatalf("MaxBy(joinPerson.ID) = %v, want %v", got, people[1])
	}
}

func TestTypedMaxByKeepsFirstEqualMaximum(t *testing.T) {
	people := []joinPerson{{id: 5, name: "first"}, {id: 5, name: "second"}}
	got := fromSlice(people).MaxBy(joinPerson.ID)

	if got != people[0] {
		t.Fatalf("MaxBy() = %v, want first maximum %v", got, people[0])
	}
}

func TestTypedMaxByReturnsZeroForEmptyQuery(t *testing.T) {
	if got := fromSlice([]int{}).MaxBy(func(value int) int { return value }); got != 0 {
		t.Fatalf("MaxBy() = %d, want zero value", got)
	}
}

func maxByBenchmarkKey(value int) int { return value }

var benchmarkMaxByResult int

func BenchmarkTypedMaxBy(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkMaxByResult = FromSlice(source).Max().(int)
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkMaxByResult = fromSlice(source).MaxBy(maxByBenchmarkKey)
		}
	})
}
