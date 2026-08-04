package linq

import "testing"

func TestTypedContainsByAcceptsMethodExpression(t *testing.T) {
	people := []joinPerson{{id: 1, name: "one"}, {id: 2, name: "two"}}
	if !fromSlice(people).ContainsBy(2, joinPerson.ID) {
		t.Fatal("ContainsBy(2, joinPerson.ID) = false, want true")
	}
}

func TestTypedContainsByReturnsFalseWithoutKey(t *testing.T) {
	if fromSlice([]int{1, 2, 3}).ContainsBy(4, func(value int) int { return value }) {
		t.Fatal("ContainsBy(4) = true, want false")
	}
}

func TestTypedContainsByStopsAtFirstMatch(t *testing.T) {
	calls := 0
	got := fromSlice([]int{1, 2, 2}).ContainsBy(2, func(value int) int {
		calls++
		return value
	})

	if !got {
		t.Fatal("ContainsBy(2) = false, want true")
	}
	if calls != 2 {
		t.Fatalf("selector called %d times, want 2", calls)
	}
}

func containsByBenchmarkKey(value int) int { return value }

var benchmarkContainsByResult bool

func BenchmarkTypedContainsBy(b *testing.B) {
	source := make([]int, 1024)
	for i := range source {
		source[i] = i
	}

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkContainsByResult = legacyFromSlice(source).Contains(1023)
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkContainsByResult = fromSlice(source).ContainsBy(1023, containsByBenchmarkKey)
		}
	})
}
