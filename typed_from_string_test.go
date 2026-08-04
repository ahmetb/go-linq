package linq

import (
	"slices"
	"strings"
	"testing"
)

func TestTypedFromStringIteratesRunes(t *testing.T) {
	type text string
	got := fromString(text("a界🙂")).toSlice()
	want := []rune{'a', '界', '🙂'}

	if !slices.Equal(got, want) {
		t.Fatalf("fromString() = %q, want %q", got, want)
	}
}

func TestTypedFromStringStopsWithConsumer(t *testing.T) {
	var got []rune
	fromString("abc").iterate(func(value rune) bool {
		got = append(got, value)
		return false
	})

	if !slices.Equal(got, []rune{'a'}) {
		t.Fatalf("fromString() yielded %q, want a", got)
	}
}

var benchmarkFromStringSum rune

func BenchmarkTypedFromString(b *testing.B) {
	source := strings.Repeat("a界🙂", 256)

	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			var sum rune
			legacyFromString(source).Iterate(func(value any) bool {
				sum += value.(rune)
				return true
			})
			benchmarkFromStringSum = sum
		}
	})

	b.Run("generic_query", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			var sum rune
			fromString(source).iterate(func(value rune) bool {
				sum += value
				return true
			})
			benchmarkFromStringSum = sum
		}
	})
}
