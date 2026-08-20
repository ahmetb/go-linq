package linq

import "testing"

const (
	size = 1000000
)

func BenchmarkSelectWhereFirst(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Range(1, size).Select(func(i int) int {
			return -i
		}).Where(func(i int) bool {
			return i > -1000
		}).First()
	}
}

func BenchmarkSelectWhereFirst_handwritten(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var r int
		for i := 1; i <= size; i++ {
			v := -i
			if v > -1000 {
				r = v
				break
			}
		}
		_ = r
	}
}

func BenchmarkSum(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Sum(Range(1, size).Where(func(i int) bool {
			return i%2 == 0
		}))
	}
}

func BenchmarkSum_handwritten(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var sum int
		for i := 1; i <= size; i++ {
			if i%2 == 0 {
				sum += i
			}
		}
		_ = sum
	}
}

func BenchmarkZipSkipTake(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Range(1, size).Zip(Range(1, size).Select(func(i int) int {
			return i * 2
		}), func(i, j int) int {
			return i + j
		}).Skip(2).Take(5).Count()
	}
}

func BenchmarkChunk(b *testing.B) {
	const chunkSize = 256
	source := make([]int, 65536)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if got := Chunk(FromSlice(source), chunkSize).Count(); got != len(source)/chunkSize {
			b.Fatalf("Chunk count=%d expected %d", got, len(source)/chunkSize)
		}
	}
}

func BenchmarkCountBy(b *testing.B) {
	const keys = 1024
	source := make([]int, 65536)
	for i := range source {
		source[i] = i % keys
	}
	q := FromSlice(source)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if got := q.CountBy(func(i int) int { return i }).Count(); got != keys {
			b.Fatalf("CountBy count=%d expected %d", got, keys)
		}
	}
}

func BenchmarkAggregateBy(b *testing.B) {
	const keys = 1024
	source := make([]int, 65536)
	for i := range source {
		source[i] = i % keys
	}
	q := FromSlice(source)
	keySelector := func(i int) int { return i }
	seedSelector := func(int) int { return 0 }
	accumulate := func(total, item int) int { return total + item }

	b.Run("Seed", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if got := q.AggregateBy(keySelector, 0, accumulate).Count(); got != keys {
				b.Fatalf("AggregateBy count=%d expected %d", got, keys)
			}
		}
	})
	b.Run("SeedSelector", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if got := q.AggregateByWithSeedSelector(keySelector, seedSelector, accumulate).Count(); got != keys {
				b.Fatalf("AggregateByWithSeedSelector count=%d expected %d", got, keys)
			}
		}
	})
}

func BenchmarkJoin(b *testing.B) {
	const keys = 1024
	source := make([]int, 65536)
	for i := range source {
		source[i] = i % keys
	}
	outer, inner := FromSlice(source), Range(0, keys)
	identity := func(i int) int { return i }
	pair := func(a, b int) int { return a + b }

	b.Run("Join", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if got := outer.Join(inner, identity, identity, pair).Count(); got != len(source) {
				b.Fatalf("Join count=%d expected %d", got, len(source))
			}
		}
	})
	b.Run("LeftJoin", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if got := outer.LeftJoin(inner, identity, identity, pair).Count(); got != len(source) {
				b.Fatalf("LeftJoin count=%d expected %d", got, len(source))
			}
		}
	})
	b.Run("RightJoin", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if got := inner.RightJoin(outer, identity, identity, pair).Count(); got != len(source) {
				b.Fatalf("RightJoin count=%d expected %d", got, len(source))
			}
		}
	})
}

func BenchmarkSequence(b *testing.B) {
	const count = 65536
	q := Sequence(0, count-1, 1)

	for n := 0; n < b.N; n++ {
		if got := Sum(q); got != count*(count-1)/2 {
			b.Fatalf("Sequence sum=%d expected %d", got, count*(count-1)/2)
		}
	}
}
