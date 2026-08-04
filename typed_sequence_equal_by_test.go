package linq

import "testing"

func TestTypedSequenceEqualByComparesDifferentTypesByMethodExpression(t *testing.T) {
	people := []joinPerson{{id: 1}, {id: 2}}
	pets := []joinPet{{ownerID: 1}, {ownerID: 2}}

	if !fromSlice(people).SequenceEqualBy(fromSlice(pets), joinPerson.ID, joinPet.OwnerID) {
		t.Fatal("SequenceEqualBy() = false, want true")
	}
}

func TestTypedSequenceEqualByDetectsLengthDifference(t *testing.T) {
	identity := func(value int) int { return value }
	if fromSlice([]int{1}).SequenceEqualBy(fromSlice([]int{1, 2}), identity, identity) {
		t.Fatal("SequenceEqualBy() = true for different lengths")
	}
}

func TestTypedSequenceEqualByStopsAtFirstMismatch(t *testing.T) {
	leftCalls, rightCalls := 0, 0
	got := fromSlice([]int{1, 2, 3}).SequenceEqualBy(
		fromSlice([]int{1, 9, 3}),
		func(value int) int { leftCalls++; return value },
		func(value int) int { rightCalls++; return value },
	)

	if got {
		t.Fatal("SequenceEqualBy() = true, want false")
	}
	if leftCalls != 2 || rightCalls != 2 {
		t.Fatalf("selectors called left=%d right=%d, want 2 each", leftCalls, rightCalls)
	}
}

func sequenceEqualBenchmarkKey(value int) int { return value }

var benchmarkSequenceEqualByResult bool

func BenchmarkTypedSequenceEqualBy(b *testing.B) {
	left := make([]int, 1024)
	right := make([]int, 1024)
	for i := range left {
		left[i] = i
		right[i] = i
	}

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkSequenceEqualByResult = FromSlice(left).SequenceEqual(FromSlice(right))
		}
	})

	b.Run("generic_method", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			benchmarkSequenceEqualByResult = fromSlice(left).SequenceEqualBy(
				fromSlice(right), sequenceEqualBenchmarkKey, sequenceEqualBenchmarkKey,
			)
		}
	})
}
