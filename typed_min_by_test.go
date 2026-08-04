package linq

import "testing"

func TestTypedMinByAcceptsMethodExpression(t *testing.T) {
	people := []joinPerson{{id: 2, name: "two"}, {id: 1, name: "one"}, {id: 3, name: "three"}}
	got := fromSlice(people).MinBy(joinPerson.ID)

	if got != people[1] {
		t.Fatalf("MinBy(joinPerson.ID) = %v, want %v", got, people[1])
	}
}

func TestTypedMinByKeepsFirstEqualMinimum(t *testing.T) {
	people := []joinPerson{{id: 1, name: "first"}, {id: 1, name: "second"}}
	got := fromSlice(people).MinBy(joinPerson.ID)

	if got != people[0] {
		t.Fatalf("MinBy() = %v, want first minimum %v", got, people[0])
	}
}

func TestTypedMinByReturnsZeroForEmptyQuery(t *testing.T) {
	if got := fromSlice([]int{}).MinBy(func(value int) int { return value }); got != 0 {
		t.Fatalf("MinBy() = %d, want zero value", got)
	}
}

func minByBenchmarkKey(value int) int { return value }

var benchmarkMinByResult int

func BenchmarkTypedMinBy(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = len(source) - i
	}

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkMinByResult = FromSlice(source).Min().(int)
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkMinByResult = fromSlice(source).MinBy(minByBenchmarkKey)
		}
	})
}
